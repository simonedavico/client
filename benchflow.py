import click
import requests
import subprocess

expManagerAddress = "http://195.176.181.55:8070"
cassandraIP = "195.176.181.55"
cassandraPORT = "9160"

@click.group()
def benchflow():
	pass

# @benchflow.command()
# def build():
# 	"""Build a BenchFlow benchmark."""
# 	cmd = ["make"]
# 	# click.echo(cmd)
# 	subprocess.call(cmd)

@benchflow.command()
def build():
	"Build a BenchFlow .zip benchmark"
	cmd = ["make", "build_for_benchflow"]
	subprocess.call(cmd)

@click.group()
def expManager():
	pass

@expManager.command()
@click.argument("benchmark", type=click.Path(exists=True), 
				metavar="<benchmark>")
def deploy(benchmark):
	"""Deploys a BenchFlow <benchmark> zip archive."""
	filename = click.format_filename(benchmark)
	file = { 'file': open(filename, 'rb') }
	click.echo("Deploying benchmark...")
	r = requests.post(expManagerAddress + "/faban/deploy", files=file)
	click.echo(r.json())

@expManager.command()
@click.argument("benchmarkshortname", metavar="<shortName>")
@click.argument("configuration", type=click.Path(exists=True),
				metavar="<configuration>")
#@click.argument("configfile", type=click.Path(exists=True))
def run(benchmarkshortname, configuration):
	"""Runs a BenchFlow benchmark 
	given its <shortName> and a <configuration>"""
	filename = click.format_filename(configuration)
	body = { 'config' : open(filename, 'rb') }
	address = expManagerAddress + "/faban/run/" + benchmarkshortname
	r = requests.post(address, files=body)
	click.echo(r.json())

@expManager.command()
@click.argument("runid", metavar="<runId>")
def status(runid):
	"""Returns the status of a benchmark run, given its <runId>"""
	r = requests.get(expManagerAddress + "/faban/status/" + runid)
	click.echo(r.json())		

@click.group()
def cassandra(args):
	pass

@cassandra.command()
# @click.argument("host", metavar="<host>")
# @click.argument("port", metavar="<port>")
def cql():
	"""Starts cqlsh on specified <host> and <port>"""
	cmd = ["cqlsh " + cassandraIP + " " + cassandraPORT]
	# click.echo(cmd)
	subprocess.call(cmd, shell=True)

benchflowClient = click.CommandCollection(sources=[benchflow, expManager, cassandra])

if __name__ == '__main__':
	benchflowClient()