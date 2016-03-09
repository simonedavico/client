import click
import requests
import subprocess
import os
import pygments
import zipfile
from pathlib import Path

exp_manager_address = os.getenv('EXPERIMENTS_MANAGER_ADDRESS')
cassandra_ip = os.getenv('CASSANDRA_IP')
cassandra_port = os.getenv('CASSANDRA_PORT')
drivers_maker_address = os.getenv('DRIVERS_MAKER_ADDRESS')


class V1(object):
    def __init__(self):
        self.session = requests.Session()
        self.session.headers = {'Accept': 'application/vnd.experiments-manager.v1+json'}

    def deploy(self, benchmark):
        filename = click.format_filename(benchmark)
        benchmark = {'benchmark': open(filename, 'rb')}
        click.echo('Deploying benchmark...')
        r = self.session.post(exp_manager_address + '/deploy', files=benchmark)
        click.echo(r.json())

    def run(self, benchmark_name, configuration):
        filename = click.format_filename(configuration)
        body = {'benchflow-benchmark': open(filename, 'rb')}
        address = '{}/run/{}'.format(exp_manager_address, benchmark_name)
        r = self.session.post(address, files=body)
        click.echo(r.json())

    def status(self, run_id):
        r = self.session.get('{}/status/{}'.format(exp_manager_address, run_id))
        click.echo(r.json())


class V2(object):
    def __init__(self):
        self.session = requests.Session()
        self.session.headers = {
            'Accept': 'application/vnd.experiments-manager.v2+json'
        }

    def deploy(self, benchmark):
        filename = click.format_filename(benchmark)
        benchmark = {'benchmark': open(filename, 'rb')}
        click.echo('Deploying benchmark...')
        r = self.session.post('{}/deploy'.format(exp_manager_address), files=benchmark)
        click.echo(r.json())

    def run(self, benchmark_name, configuration):
        filename = click.format_filename(configuration)
        # TODO: enable config override
        # body = {'benchflow-benchmark': open(filename, 'rb')}
        address = '{}/run/{}'.format(exp_manager_address, benchmark_name)
        r = self.session.post(address)
        click.echo(r.json())

    def status(self, run_id):
        raise click.ClickException('/status is not implemented for this api version yet.')


class Config(object):
    def __init__(self):
        self.version = 2
        self.api = V2()


pass_config = click.make_pass_decorator(Config, ensure=True)


@click.group()
@click.option('--version', default=2, help='api version')
def cli():
    pass


def zipdir(path):
    """Utility function to zip a directory"""
    p = Path(path)
    parent = (p / '..').resolve()
    archive_path = '{}/{}.zip'.format(parent, p.name)
    with zipfile.ZipFile(archive_path, 'w', zipfile.ZIP_DEFLATED) as archive:
        for root, dirs, files in os.walk(path):
            for file in files:
                archive.write(os.path.join(root, file),
                              os.path.relpath(os.path.join(root, file), os.path.join(path, '..')))
    return archive_path


@cli.command()
@click.argument('benchmark_dir', type=click.Path(exists=True, file_okay=False),
                metavar='<benchmark_dir>')
@pass_config
def build(config, benchmark_dir):
    """Builds a BenchFlow benchmark"""
    if config.version == 'v2':
        p = Path(benchmark_dir)
        dd = (p / 'docker-compose.yml').resolve()
        bb = (p / 'benchflow-benchmark.yml').resolve()
        models = (p / 'models').resolve()
        sources = (p / 'sources').resolve()
        sources_zip_path = zipdir(str(sources))
        benchmark_zip_path = '{}/{}.zip'.format(benchmark_dir, p.name)
        with zipfile.ZipFile(benchmark_zip_path, 'w', zipfile.ZIP_DEFLATED) as archive:
            archive.write(str(dd), dd.name)
            archive.write(str(bb), bb.name)
            archive.write(str(models), models.name)
            archive.write(sources_zip_path, 'sources.zip')
        os.remove(sources_zip_path)
        benchmark_name = click.style(p.name, fg='red')
        click.echo('Benchmark {} successfully built.'.format(benchmark_name))
    else:
        cmd = ['make', 'build_for_benchflow']
        subprocess.call(cmd)


@cli.command()
def update():
    """Updates BenchFlow by pulling the latest image"""
    pass


@click.group()
def api():
    pass


@api.command()
@click.argument('benchmark', type=click.Path(exists=True, dir_okay=False),
                metavar='<benchmark>')
@pass_config
def deploy(config, benchmark):
    """Deploys a <benchmark>"""
    config.api.deploy(benchmark)


@api.command()
@click.argument('benchmark_name', metavar='<benchmark_name>')
@click.argument('configuration', type=click.Path(exists=True, dir_okay=False),
                metavar='<benchmark_configuration>', required=False)
@pass_config
def run(config, benchmark_name, configuration):
    """Runs a benchmark"""
    config.api.run(benchmark_name, configuration)


@api.command()
@click.argument('run_id', metavar='<run_id>')
@pass_config
def status(config, run_id):
    """Returns the status of a benchmark run"""
    config.api.status(run_id)


@cli.command()
def cql():
    """Starts cqlsh"""
    cmd = ['cqlsh ' + cassandra_ip + ' ' + cassandra_port]
    subprocess.call(cmd, shell=True)


@cli.group()
def debug():
    """Debug commands, for development purposes"""
    pass


@debug.command()
@click.argument('configuration', type=click.Path(exists=True, dir_okay=False),
                metavar='<configuration>')
def convert(configuration):
    """Generates a Faban configuration file"""
    filename = click.format_filename(configuration)
    bfconfiguration = {'benchflow-benchmark': open(filename, 'rb')}
    click.echo('Address: ' + drivers_maker_address + '/convert')
    r = requests.post(drivers_maker_address + '/convert', files=bfconfiguration)
    click.echo(r.text)


@click.command(cls=click.CommandCollection, sources=[api, cli])
@click.option('--api-version', default='v2', type=click.Choice(['v1', 'v2']), help='Api version')
@pass_config
def client(config, api_version):
    config.version = api_version
    config.api = V1() if config.version == 'v1' else V2()
