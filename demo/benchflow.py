import click
import requests
import subprocess

expManagerAddress = "http://localhost:8080"

@click.group()
def expManager():
	pass	

@expManager.command()
@click.argument("benchmark", type=click.Path(exists=True))
def deploy(benchmark):
	"""Deploys a BenchFlow benchmark"""
	filename = click.format_filename(benchmark)
	file = { 'file': open(filename, 'rb') }
	r = requests.post(expManagerAddress + "/faban/deploy", files=file)
	click.echo(r.json())

@expManager.command()
@click.argument("benchmarkshortname")
#@click.argument("configfile", type=click.Path(exists=True))
def run(benchmarkshortname):
	"""Runs a BenchFlow benchmark"""
	r = requests.post(expManagerAddress + "/faban/run/" + benchmarkshortname)
	click.echo(r.json())

@expManager.command()
@click.argument("runid")
def status(runid):
	"""Returns the status of a benchmark run"""
	r = requests.get(expManagerAddress + "/faban/status/" + runid)
	click.echo(r.json())		

@click.group()
def cassandra(args):
	pass

@cassandra.command()
@click.argument("host")
@click.argument("port")
def clsqh(host, port):
	"""Starts clsqh"""
	cmd = list(("clsqh", host, port))
	click.echo(cmd)
	subprocess.run(cmd)

benchflowClient = click.CommandCollection(sources=[expManager, cassandra])

if __name__ == '__main__':
	benchflowClient()