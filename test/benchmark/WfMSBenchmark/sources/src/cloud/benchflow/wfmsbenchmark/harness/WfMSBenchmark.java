package cloud.benchflow.wfmsbenchmark.harness;

//TODO: remove not useful imports, and add the ones that should be here even if it is not needed (e.g., String)

import com.sun.faban.harness.Validate;
import com.sun.faban.harness.Configure;
import com.sun.faban.harness.PreRun;
import com.sun.faban.harness.StartRun;
import com.sun.faban.harness.PostRun;
import com.sun.faban.harness.EndRun;
import com.sun.faban.harness.KillRun;
import com.sun.faban.harness.ParamRepository;
import com.sun.faban.harness.DefaultFabanBenchmark2;

import com.sun.faban.harness.RunContext;

import org.w3c.dom.Element;

import java.io.*;
import java.util.HashMap;
import java.util.Map;
import java.util.logging.Logger;


// //WfMS interaction specific
import com.google.gson.JsonArray;
import com.google.gson.JsonObject;
import com.google.gson.JsonParser;

import com.sun.faban.driver.transport.hc3.ApacheHC3Transport;

import org.apache.commons.httpclient.methods.PutMethod;
import org.apache.commons.httpclient.methods.multipart.FilePart;
import org.apache.commons.httpclient.methods.multipart.MultipartRequestEntity;
import org.apache.commons.httpclient.methods.multipart.Part;
import org.apache.commons.httpclient.methods.multipart.StringPart;

import org.w3c.dom.Document;
import org.w3c.dom.Element;
import org.w3c.dom.Node;
import org.w3c.dom.NodeList;

import java.util.List;
import java.util.ArrayList;
import java.util.LinkedList;

import java.util.concurrent.*;
// import java.util.concurrent.Callable;
// import java.util.concurrent.ExecutionException;
// import java.util.concurrent.ExecutorService;
// import java.util.concurrent.Executors;
// import java.util.concurrent.Future;


/**
 *
 * @author Vincenzo Ferme
 */
public class WfMSBenchmark extends DefaultFabanBenchmark2 {

    private static Logger logger = Logger.getLogger(cloud.benchflow.wfmsbenchmark.harness.WfMSBenchmark.class.getName());

    private static String benchmarkDir;

    private String SUTEndpoint;

    private ApacheHC3Transport http;

    private JsonParser parser;

    private WfMSBenchmarkApi wfms;

    protected ParamRepository params;

    public abstract class WfMSBenchmarkApi {

        protected String sutEndpoint;
        protected String deployAPI;

        public WfMSBenchmarkApi(String sutEndpoint, String deployAPI) {
            this.sutEndpoint = sutEndpoint;
            this.deployAPI = sutEndpoint + deployAPI;
        };

        public abstract Map<String, String> deploy(File model) throws IOException;

    }

    private class WfMSBenchmarkImpl extends WfMSBenchmarkApi {

        protected String processDefinitionAPI;

        public WfMSBenchmarkImpl(String sutEndpoint) {
            super(sutEndpoint, "/deployment/create");
            this.processDefinitionAPI = sutEndpoint + "/process-definition";
        }

        @Override
        public Map<String, String> deploy(File model) throws IOException {

            Map<String, String> result = new HashMap<String, String>();
            StringPart deploymentName = new StringPart("deployment-name", model.getName());
            logger.info("Deploying model: " + model.getAbsolutePath());
            List<Part> parts = new ArrayList<Part>();

            FilePart process = new FilePart("*", model);

            parts.add(deploymentName);
            parts.add(process);
            logger.info("Deploying model at: " + deployAPI);
            StringBuilder deployDef = http.fetchURL(deployAPI, parts);

            logger.info("DEPLOYMENT RESPONSE: " + deployDef.toString());
            JsonObject deployObj = parser.parse(deployDef.toString()).getAsJsonObject();
            String deploymentId = deployObj.get("id").getAsString();

            logger.info("DEPLOYMENT ID: " + deploymentId);

            //Obtain process definition data
            StringBuilder procDef = http.fetchURL(processDefinitionAPI + "?deploymentId=" + deploymentId);
            logger.info("PROCESS DEFINITION RESPONSE: " + procDef.toString());
            String processDefinitionResponse = procDef.toString();

            JsonArray procDefArray = parser.parse(processDefinitionResponse).getAsJsonArray();
            //We only get 1 element using the deploymentId
            String processDefinitionId = procDefArray.get(0).getAsJsonObject().get("id").getAsString();
            result.put(model.getName(), processDefinitionId);
            return result;

        }
    }


    /**
     * TODO: Validate the configuration
     * Preprocess the configuration so that it is more easly accessile by the drivers
     */
     @Override
     @Validate public void validate() throws Exception {

        super.validate();

        logger.info("START: Validate...");   

        initialize();

        logger.info("DONE: initialize");   

        logger.info("benchmarkDir is: " + benchmarkDir);   

        setSUTEndpoint();

        wfms = new WfMSBenchmarkImpl(SUTEndpoint);

        logger.info("DONE: setSUTEndpoint");   

        // addConfigurationNode("sutConfig","SUTEndpoint",SUTEndpoint);
        addConfigurationNode("sutConfiguration","SUTEndpoint",SUTEndpoint);

        logger.info("DONE: addConfigurationNode");  

        logger.info("END: Validate...");

     }


     /**
      * Deploy the SUT, monitors and collectors, Start it and Check that everything is UP and Running
      *
      * @throws Exception If configuration was not successful
      */
     @Configure public void configure() throws Exception {
          //simone
    	  //String benchFlowComposeService = getXPathValue("services/benchFlowCompose");
          String benchFlowComposeService = getXPathValue("benchFlowServices/benchFlowCompose"); 
          //simone
    	  //String runID = getXPathValue("runInfo/runID");
          String runID = getXPathValue("benchFlowRunConfiguration/trialId");    	 
    	  
          String sutDir = benchmarkDir + "sut";
    	  
          //simone
    	  //File dockerCompose = new File(sutDir + "/docker-compose.yml");
            File dockerCompose = new File(sutDir + "/docker-compose-" + runID + ".yml");
 		 
     		//Deploying the system under test
     		//curl -v -X PUT -F 'docker_compose_file=@docker-compose.yml' 
     		//-F 'benchflow_compose_file=@docker-compose.yml' 
     		//http://<HOST_IP>:<HOST_PORT>/projects/camunda/deploymentDescriptor/		      
        FilePart dockerComposeFile = new FilePart("docker_compose_file", dockerCompose);
         
//         URL deployAPI = new URL(endPoint + "/projects/camunda/deploymentDescriptor/");
        String deployAPI = benchFlowComposeService + "/projects/" + runID + "/deploymentDescriptor/";
 		
        PutMethod put = new PutMethod(deployAPI);
        
        Part[] partsArray = {
        		dockerComposeFile
        };
        
        put.setRequestEntity(
            new MultipartRequestEntity(partsArray, put.getParams())
        );

        int status = http.getHttpClient().executeMethod(put);
        
        logger.info("System Deployed. Status: " + status);
        
        //start the system under test
//      curl -v -X PUT http://<HOST_IP>:<HOST_PORT>/projects/camunda/up/
        String upAPI = benchFlowComposeService + "/projects/" + runID + "/up/";
      		
        PutMethod putUp = new PutMethod(upAPI);
      
        int statusUp = http.getHttpClient().executeMethod(putUp);
      
        logger.info("System Started. Status: " + statusUp);

     }


    /**
     * Deploy the models and update the run.xml file with the data about them
     * TODO: of course, for everything, error handling in case of problem interacting with the SUT etc..
     *
     * @throws Exception If configuration was not successful
     */
    @PreRun public void preRun() throws Exception {

        //The following code is custom for Camunda, but all the engines exposing rest APIs should work similarly

        logger.info("START: Deployng processes...");

        int numDeplProcesses = 0;

        String modelDir = benchmarkDir + "models";

        String deployAPI = SUTEndpoint + "/deployment/create";

        String processDefinitionAPI = SUTEndpoint + "/process-definition";

        File folder = new File(modelDir);

        File[] listOfFiles = folder.listFiles();

        //Add models node
        String agentName = "WfMSBenchmarkDriver";
        String driverToUpdate = "fa:runConfig/fd:driverConfig[@name=\"" + agentName + "\"]";
        //Here I am assuming there is not an already defined properties element
        Element properties = addConfigurationNode(driverToUpdate,"properties","");

        for (int i = 0; i < listOfFiles.length; i++) {
            if (listOfFiles[i].isFile()) {

                String modelName = listOfFiles[i].getName();
                String modelPath = modelDir + "/" + modelName;

                File modelFile = new File(modelPath);

                String processDefinitionId = wfms.deploy(modelFile).get(modelName);
                logger.info("PROCESS DEFINITION ID: " + processDefinitionId);
                addModel(properties, i+1, modelName,processDefinitionId);
                numDeplProcesses++;
            
            }
        }

        addProperty(properties, "model_num", String.valueOf (numDeplProcesses));
        
        logger.info("END: Deployng processes...");

    }

    //NOTE: it was in the OnceBefore method of the drivers, but this point it is more suitable,
    //because it is just when the ramp up starts
    /**
     * Start monitors and collectors
     * Only thread 0 does is and can share data with onceAfter
     */
    @Override
    @StartRun public void start() throws Exception {

        super.start();

        // You will need to ensure all the driver processes on all driver systems get started, and - if feasible - enter the rampup state before returning from this method as the tools timer will get started immediately after this method terminates.
        //This means that starting the stats here might miss some data point at the beginning of the run, but
        //we anyway delete some of the process at the beginning, part of the warm up period, so it is fine


        // Stats:
        // TODO: the following calls should be done in parallel
         // creating thread pool to execute task which implements Callable
        // The best point to start the stats collection should be inside the start method,
        // just before pinging the master to start the RAMP UP

        //TODO: this is suboptimal because we create one thread for each collector
        //but we should really create one thread only for the collectors with a start API
        //this will allow us also to count the number of expected responses and use a countdownlatch 
        NodeList collectors = getNodes("benchFlowServices/collectors"); 
        ExecutorService es = Executors.newFixedThreadPool(collectors.getLength());
        CompletionService<String> cs = new ExecutorCompletionService<>(es);
        // - curl -v -X GET http://<HOST_IP>:<HOST_PORT>/start
        
        // String statsDBStart = getXPathValue("services/collectors/statsDBStart");
        // logger.info("STATS DB START URL: " + statsDBStart);
        // // StringBuilder resDB = http.fetchURL(statsDBStart);
        // Future<String> resDBFuture = es.submit(new BenchFlowServicesAsynchInteraction(statsDBStart));
        // //simone
        // String statsWfMSStart = getXPathValue("services/collectors/statsWfMSStart");
        // logger.info("STATS WfMS START URL: " + statsWfMSStart);
        // // StringBuilder resWfMS = http.fetchURL(statsWfMSStart);
        // Future<String> resWfMSFuture = es.submit(new BenchFlowServicesAsynchInteraction(statsWfMSStart));   
        
        // String resDB = resDBFuture.get();
        // logger.info("STATS DB START: " + resDB);
        // String resWfMS = resWfMSFuture.get();
        // logger.info("STATS WfMS START: " + resWfMS);

        //simone
        List<Future<String>> collectorsStartResponses = new LinkedList<Future<String>>();
        for(int i = 0; i < collectors.getLength(); i++) {
            Node collector = collectors.item(i);
            NodeList collectorAPIs = collector.getChildNodes(); 
            for(int j = 0; j < collectorAPIs.getLength(); j++) {
                Node collectorAPI = collectorAPIs.item(j);
                if(collectorAPI.getNodeName().equals("start")) {
                    String collectorStartAPI = collectorAPI.getNodeValue();
                    logger.info("Starting service " + collectorAPI.getNodeName());
                    //TODOsimone: update collectorAPI with complete address
                    collectorsStartResponses.add(cs.submit(new BenchFlowServicesAsynchInteraction(collectorStartAPI)));
                }
            }
        }

        for(Future<String> collectorStartResponse: collectorsStartResponses) {
            String response = collectorStartResponse.get();
            logger.info("Start response: " + collectorStartResponse);
        }
    }


    // /**
    //  * Collect data and undeploy the SUT, monitors and collectors checking that they are correclty undeployed
    //  *
    //  * @throws Exception If configuration was not successful
    //  */
    @PostRun public void postRun() throws Exception {
    
        //collect the data
        // MysqlDump:
        // - curl -v -X GET http://<HOST_IP>:<HOST_PORT>/data
        //simone
        // String mysqlDumpGetData = getXPathValue("services/collectors/mysqldumpData");
        // http.fetchURL(mysqlDumpGetData);

        //TODO: this is suboptimal because we create one thread for each collector
        //but we should really create one thread only for the collectors with a start API
        //this will allow us also to count the number of expected responses and use a countdownlatch 
        NodeList collectors = getNodes("benchFlowServices/collectors"); 
        ExecutorService es = Executors.newFixedThreadPool(collectors.getLength());
        CompletionService<String> cs = new ExecutorCompletionService<>(es);

        // Stats:
        // TODO: the following calls should be done in parallel
        // - curl -v -X GET http://<HOST_IP>:<HOST_PORT>/stop
        //simone
        // String statsDBStop = getXPathValue("services/collectors/statsDBStop");
        // logger.info("STATS DB STOP URL: " + statsDBStop);
        // // StringBuilder resDB = http.fetchURL(statsDBStop);
        // Future<String> resDBFuture = es.submit(new BenchFlowServicesAsynchInteraction(statsDBStop));

        // //simone
        // String statsWfMSStop = getXPathValue("services/collectors/statsWfMSStop");
        // logger.info("STATS WfMS STOP URL: " + statsWfMSStop);
        // // StringBuilder resWfMS = http.fetchURL(statsWfMSStop);
        // Future<String> resWfMSFuture = es.submit(new BenchFlowServicesAsynchInteraction(statsWfMSStop));

        // String resDB = resDBFuture.get();
        // logger.info("STATS DB STOP: " + resDB);
        // String resWfMS = resWfMSFuture.get();
        // logger.info("STATS WfMS STOP: " + resWfMS);

        List<Future<String>> collectorsStopResponses = new LinkedList<Future<String>>();
        for(int i = 0; i < collectors.getLength(); i++) {
            Node collector = collectors.item(i);
            NodeList collectorAPIs = collector.getChildNodes(); 
            for(int j = 0; j < collectorAPIs.getLength(); j++) {
                Node collectorAPI = collectorAPIs.item(j);
                if(collectorAPI.getNodeName().equals("stop")) {
                    String collectorStopAPI = collectorAPI.getNodeValue();
                    logger.info("Stopping service " + collectorAPI.getNodeName());
                    //TODOsimone: update collectorAPI with complete address
                    collectorsStopResponses.add(cs.submit(new BenchFlowServicesAsynchInteraction(collectorStopAPI))); 
                }
            }
        }

        for(Future<String> collectorStopResponse: collectorsStopResponses) {
            String response = collectorStopResponse.get();
            logger.info("Stop response: " + collectorStopResponse);
        }        
    	 
    	//remove the sut
        //curl -v -X PUT http://<HOST_IP>:<HOST_PORT>/projects/camunda/rm/
    	//simone
    	// String benchFlowComposeService = getXPathValue("services/benchFlowCompose");
        String benchFlowComposeService = getXPathValue("benchFlowServices/benchFlowCompose");
    	// String runID = getXPathValue("runInfo/runID");
        String runID = getXPathValue("benchFlowRunConfiguration/trialId");
        String rmAPI = benchFlowComposeService + "/projects/" + runID + "/rm/";
        PutMethod putRm = new PutMethod(rmAPI);
        int statusRm = http.getHttpClient().executeMethod(putRm);
        logger.info("System UnDeployed. Status: " + statusRm);
     }

     @Override
     @EndRun public void end() throws Exception {
         try {
            super.end(); 
         } catch (Exception e) {
            //remove the sut
            //curl -v -X PUT http://<HOST_IP>:<HOST_PORT>/projects/camunda/rm/
            //simone
            // String benchFlowComposeService = getXPathValue("services/benchFlowCompose");
            String benchFlowComposeService = getXPathValue("benchFlowServices/benchFlowCompose");
            // String runID = getXPathValue("runInfo/runID");
            String runID = getXPathValue("benchFlowRunConfiguration/trialId");
            String rmAPI = benchFlowComposeService + "/projects/" + runID + "/rm/";
            PutMethod putRm = new PutMethod(rmAPI);
            int statusRm = http.getHttpClient().executeMethod(putRm);
            logger.info("System UnDeployed in End. Status: " + statusRm);
         }
     }

    /***
    * Undeploy the SUT, monitors and collectors checking that they are correclty undeployed
    *
    * This is called in case of FAILED execution, instead of the PostRun [TODO: verify that this is the case, or we ask twice for undeploy]
    * @throws Exception If configuration was not successful
    */
    @KillRun public void kill() throws Exception {
        //remove the sut
        //curl -v -X PUT http://<HOST_IP>:<HOST_PORT>/projects/camunda/rm/
        //simone
        // String benchFlowComposeService = getXPathValue("services/benchFlowCompose");
        String benchFlowComposeService = getXPathValue("benchFlowServices/benchFlowCompose");
        // String runID = getXPathValue("runInfo/runID");
        String runID = getXPathValue("benchFlowRunConfiguration/trialId");
        String rmAPI = benchFlowComposeService + "/projects/" + runID + "/rm/";
        PutMethod putRm = new PutMethod(rmAPI);
        int statusRm = http.getHttpClient().executeMethod(putRm);
        logger.info("System UnDeployed. Status: " + statusRm); 
    }

    /***
    * Things that can be abstracted away
    * Some of them must also be shared with the Drivers, since they are basically duplicated code
    */
    private void initialize() {
        params = RunContext.getParamRepository(); 
        benchmarkDir = RunContext.getBenchmarkDir();
        http = new ApacheHC3Transport();
        parser = new JsonParser();
    }

    private void addModel(Element properties, int modelNum, String modelName, String processDefinitionId) throws Exception {
        //We need to attach them as driver properties otherwise it is not possible to access them in the Driver
        //Add the information about the deployed process in the run context
        //TODO: provide abstracted method to improve the adding of informations like the following, dinamically
        //Maybe also improving com.sun.faban.harness.ParamRepository if needed 
        /**
          * <models>
          *  <model id="processDefinitionId">
          *   <name></name>
          *   <startID></startID>
          *  </model>
          * </models>
          */

          // <property name="path3">spidermark/</property>
          addProperty(properties, "model_" + modelNum + "_name", modelName);
          addProperty(properties, "model_" + modelNum + "_startID", processDefinitionId);
    }

    private void addProperty(Element properties, String name, String value) throws Exception {

        //TODO, move and decide where to put this code, because basically I'm getting the run document and working directly on it
        //maybe we should move this stuff is com.sun.faban.harness.ParamRepository
        // root elements
        //simone TODO: check if this way of retrieving the document works
        // Document runDoc = params.getNode("benchFlowBenchmark").getOwnerDocument();
        Document runDoc = params.getTopLevelElements().item(0).getOwnerDocument();
        Element prop = addConfigurationNode(properties,"property","");
        prop.setAttribute("name",name);
        prop.appendChild(runDoc.createTextNode(value));
        properties.appendChild(prop);
        params.save();
    }

    //TODO: avoid multiple params.save(); calls and setup a dedicated method to be called at the end of each Faban Driver operations

    private Element addConfigurationNode(String baseXPath, String nodeName, String value) throws Exception {
        Element node = params.addParameter(baseXPath, null, null, nodeName);
        params.setParameter(node, value);
        params.save();
        return node;
    }

    private Element addConfigurationNode(Element parent, String nodeName, String value) throws Exception {
        Element node = params.addParameter(parent, null, null, nodeName);
        params.setParameter(node, value);
        params.save();
        return node;
    }

    private String getXPathValue(String xPathExpression) throws Exception {
        return params.getParameter(xPathExpression);
    }


    private NodeList getNodes(String xPathExpression) {
        return params.getNodes(xPathExpression);
    }

    private void setSUTEndpoint() throws Exception {
        StringBuilder urlBuilder = new StringBuilder();
        // //Set the protocol
        // String s = getXPathValue("sutConfig/protocol");
        // // logger.info("protocol: " + s);
        // urlBuilder.append(s).append("://");
        // s = getXPathValue("sutConfig/hostPorts");
        // // logger.info("hostPorts: " + s);
        // urlBuilder.append(s);
        // //Set the contextPath
        // s = getXPathValue("sutConfig/contextPath");
        // // logger.info("contextPath: " + s);
        // if (s.charAt(0) == '/')
        //     urlBuilder.append(s);
        // else
        //     urlBuilder.append('/').append(s);

        urlBuilder.append("http://")
                  .append(getXPathValue("sutConfiguration/serviceAddress"))
                  .append(getXPathValue("sutConfiguration/endpoint"));

        SUTEndpoint = urlBuilder.toString();
    }

    class BenchFlowServicesAsynchInteraction implements Callable<String> {
        private String url;
        
        public BenchFlowServicesAsynchInteraction(String url){
            this.url = url;
        }

        @Override
        public String call() throws Exception {
            return http.fetchURL(url).toString();
        }

    }

}
