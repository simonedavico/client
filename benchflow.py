import click
import requests
import subprocess
import os
import json
from pygments import highlight
from pygments.formatters import Terminal256Formatter as TermF
from pygments.lexers import JsonLexer as JsonL


expManagerAddress = os.getenv("EXPERIMENTS_MANAGER_ADDRESS")
cassandraIP = os.getenv("CASSANDRA_IP")
cassandraPORT = os.getenv("CASSANDRA_PORT")
driversMakerAddress = os.getenv("DRIVERS_MAKER_ADDRESS")

class client:
	"""Contains utility functions for the BenchFlow client"""
	@staticmethod
	def echo(json, colors=True):
		"""Echoes highlighted json. Generalise to handle every kind of input?"""
		def color():
			return highlight(json.dumps(json), JsonL(), TermF(style='fruity'))
		click.echo(color() if colors else  json)

@click.group()
def benchflow():
	pass

# @benchflow.command()
# def foobar():
# 	client.echo('{ "foo" : "bar" }', colors=True)

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
	file = { 'benchmark': open(filename, 'rb') }
	click.echo("Deploying benchmark...")
	r = requests.post(expManagerAddress + "/deploy", files=file) 
	click.echo(r.json())
	# client.echo(r.json(), colors=True)

@expManager.command()
@click.argument("benchmarkshortname", metavar="<shortName>")
@click.argument("configuration", type=click.Path(exists=True),
				metavar="<configuration>")
#@click.argument("configfile", type=click.Path(exists=True))
def run(benchmarkshortname, configuration):
	"""Runs a BenchFlow benchmark 
	given its <shortName> and a <configuration>"""
	filename = click.format_filename(configuration)
	body = { 'benchflow-benchmark' : open(filename, 'rb') }
	address = expManagerAddress + "/run/" + benchmarkshortname
	r = requests.post(address, files=body)
	# click.echo(r.text)
	client.echo(r.json())

@expManager.command()
@click.argument("runid", metavar="<runId>")
def status(runid):
	"""Returns the status of a benchmark run, given its <runId>"""
	r = requests.get(expManagerAddress + "/status/" + runid)
	client.echo(r.json())		

@click.group()
def cassandra(args):
	pass

@cassandra.command()
def cql():
	"""Starts cqlsh on specified <host> and <port>"""
	cmd = ["cqlsh " + cassandraIP + " " + cassandraPORT]
	subprocess.call(cmd, shell=True)

@click.group()
def debug():
	pass

@debug.command()
@click.argument("configuration", type=click.Path(exists=True),
				metavar="<configuration>")
def convert(configuration):
	filename = click.format_filename(configuration)
	bfconfiguration = { 'benchflow-benchmark': open(filename, 'rb') }
	click.echo('Address: ' + driversMakerAddress + '/convert')
	r = requests.post(driversMakerAddress + '/convert', files=bfconfiguration)
	click.echo(r.text)


benchflowClient = click.CommandCollection(sources=[benchflow, expManager, cassandra, debug])

if __name__ == '__main__':	
	benchflowClient()