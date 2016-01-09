import click
import requests

expManager = "localhost:9980"

@click.group()
def faban():
	pass

@faban.command()
@click.argument("benchmark", type=click.Path(exists=True))
def deploy(benchmark):
	filename = click.format_filename(benchmark)
	file = { 'file': open(filename, 'rb') }
	r = requests.post(expManager + "/faban/deploy", files=file)
	click.echo(r.json())

@faban.command()
@click.argument("benchmarkshortname")
#@click.argument("configfile", type=click.Path(exists=True))
def run(benchmarkshortname):
	r = requests.post(expManager + "/faban/run/" + benchmarkshortname)
	click.echo(r.json())

@faban.command()
@click.argument("runid")
def status(runid):
	r = requests.get(expManager + "/faban/status/" + runid)
	click.echo(r.json())		

if __name__ == '__main__':
	faban()