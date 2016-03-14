package cloud.benchflow.wfmsbenchmark.driver;

//TODO: remove not useful imports, and add the ones that should be here even if it is not needed (e.g., String)

import com.sun.faban.common.NameValuePair;
import com.sun.faban.common.Utilities;
import com.sun.faban.driver.*;
import com.sun.faban.driver.util.ContentSizeStats;

import javax.xml.xpath.XPathExpressionException;
import javax.naming.ConfigurationException;
import java.io.IOException;
import java.util.List;
import java.util.concurrent.TimeUnit;
//Added for debugging and logging from the driver to the Harness
import java.util.logging.Logger;

//WfMS interaction specific 
import com.google.gson.Gson;

import org.w3c.dom.NodeList;
import org.w3c.dom.Node;

import org.apache.commons.httpclient.methods.multipart.FilePart;
import org.apache.commons.httpclient.methods.multipart.Part;
import org.apache.commons.httpclient.methods.multipart.StringPart;

import java.io.File;
import java.io.IOException;
import java.util.ArrayList;
import java.util.HashMap;
import java.util.Map;
import java.util.TreeMap;

import java.util.concurrent.Callable;
import java.util.concurrent.ExecutionException;
import java.util.concurrent.CountDownLatch;

/***
 * Basic WfMS workload driver drives only one operation via a get request.
 *
 * @author Vincenzo Ferme
 */
@BenchmarkDefinition (
    name    = "WfMSBenchmark Workload",
    version = "0.1",
    metric  = "req/s"
)
@BenchmarkDriver (
    name             = "WfMSBenchmarkDriver",
    threadPerScale   = (float)1,
    opsUnit          = "requests",
    metric           = "req/s",
    responseTimeUnit = TimeUnit.MILLISECONDS//,
    // percentiles = { "90", "95th", "99.9th%"} // Show different supported formats
)

@MatrixMix (
    operations = {"activitiModel", "camundaAdditionalApproval", "camundaCounterpartyOnboarding", "camundaInvoiceCollaboration", "camundaOop"},
    mix = { @Row({ 50, 50, 50, 50 ,50 }),
            @Row({ 50, 50, 50, 50 ,50 }),
            @Row({ 50, 50, 50, 50 ,50 }),
            @Row({ 50, 50, 50, 50 ,50 }),
            @Row({ 50, 50, 50, 50 ,50 }) },
    deviation = 5
)
// @MatrixMix (
//     operations = {"activitiModel", "camundaAdditionalApproval", "camundaInvoiceCollaboration", "camundaOop"},
//     mix = { @Row({ 50, 50, 50, 50 }),
//             @Row({ 50, 50, 50, 50 }),
//             @Row({ 50, 50, 50, 50 }),
//             @Row({ 50, 50, 50, 50 }) },
//     deviation = 5
// )

@FixedTime (
    cycleType = CycleType.THINKTIME,
    cycleTime = 1000,
    cycleDeviation = 5
)
public class WfMSBenchmarkDriver {

    private DriverContext ctx;
    private HttpTransport http;
    private String SUTEndpoint;
    // ContentSizeStats contentStats = null;
    private Map<String, String> JSONHeaders = new TreeMap<String, String>();
    private Map<String,String> modelsStartID = new HashMap<String,String>();
    private WfMSBenchmarkDriverApi wfms;
    //Added for debugging and logging from the driver to the Harness
    Logger logger;

    public abstract class WfMSBenchmarkDriverApi {

        protected String sutEndpoint;

        public WfMSBenchmarkDriverApi(String sutEndpoint) {
            this.sutEndpoint = sutEndpoint;
        }

        public abstract String startProcessDefinition(String processDefinitionId) throws IOException;

    }

    private class WfMSBenchmarkDriverImpl extends WfMSBenchmarkDriverApi {


        public WfMSBenchmarkDriverImpl(String sutEndpoint) {
            super(sutEndpoint);
        }

        @Override
        public String startProcessDefinition(String modelName) throws IOException {
            String startURL = SUTEndpoint + "/process-definition/" + modelsStartID.get(modelName) + "/start";
            StringBuilder responseStart = http.fetchURL(startURL, "{}", JSONHeaders);
            return responseStart.toString();
        }
    }


    /***
     * Constructs the basic web workload driver.
     * @throws XPathExpressionException An XPath error occurred
     */
    public WfMSBenchmarkDriver() throws XPathExpressionException, ConfigurationException {
        initialize();
        setSUTEndpoint();
        loadModelsInfo();
        this.wfms = new WfMSBenchmarkDriverImpl(SUTEndpoint);
    }

    //TODO: the following two methods can be abstracted away
    /***
    * Start monitors and collectors
    * Only thread 0 does is and can share data with onceAfter
    */
    @OnceBefore
    public void onceBefore() throws XPathExpressionException, InterruptedException, ExecutionException, IOException {
        //We wait a bit to create a gap in the data (TODO-RM: experimenting with data cleaning)
        //and be sure the model started during the warm up and timing synch of the sistem, end,
        //event though now that we use mock models they end very fast
        try {
            Thread.sleep(20000);
        } catch (InterruptedException e) {
        }
        logger.info("Tested pre-run (sleep 20) done");
        // An @OnceBefore operation gets called by global thread 0 of the driver only once and before any other driver thread gets instantiated. This will guarantee that no operation is called while the  @OnceBefore operation is running.
        // This also means that starting stats here collect a lot of points during the initialziation of other agetns
    //     // Stats:
    //     // TODO: the following calls should be done in parallel
    //      // creating thread pool to execute task which implements Callable
    //     ExecutorService es = Executors.newSingleThreadExecutor();
    //     // - curl -v -X GET http://<HOST_IP>:<HOST_PORT>/start
    //     String statsDBStart = getXPathValue("services/collectors/statsDBStart");
    //     logger.info("STATS DB START URL: " + statsDBStart);
    //     // StringBuilder resDB = http.fetchURL(statsDBStart);
    //     Future<String> resDBFuture = es.submit(new BenchFlowServicesAsynchInteraction(statsDBStart));

    //     String statsWfMSStart = getXPathValue("services/collectors/statsWfMSStart");
    //     logger.info("STATS WfMS START URL: " + statsWfMSStart);
    //     // StringBuilder resWfMS = http.fetchURL(statsWfMSStart);
    //     Future<String> resWfMSFuture = es.submit(new BenchFlowServicesAsynchInteraction(statsWfMSStart));   
        

    //     String resDB = resDBFuture.get();
    //     logger.info("STATS DB START: " + resDB);
    //     String resWfMS = resWfMSFuture.get();
    //     logger.info("STATS WfMS START: " + resWfMS);
    }

    /***
     * Ask monitors to wait for completion
     * Only thread 0 does is and can share data with onceBefore
     */
    @OnceAfter
    public void onceAfter() {

        // Currently we just use the MySQL monitor, since the overhead is minimal
        // curl "http://192.168.41.105:9303/status?query=SELECT+COUNT(*)+FROM+ACT_HI_PROCINST+WHERE+END_TIME_+IS+NULL&value=0&method=equal"
        //TODO: check what actually happens if one of the following exception is thrown
        final CountDownLatch done = new CountDownLatch(1);
        
        new Thread(new Runnable() {
            @Override
            public void run() {
                //TODO: use setParams
                String mysqlMonitorEndpoint = "";
                //TODO: we for sure want to have a better way to get the same.
                //the point now is that it is not possible to throw an exception on the run method
                try {
                    //simone: change this
                    // mysqlMonitorEndpoint = getXPathValue("services/monitors/mysql");
                    mysqlMonitorEndpoint = getXPathValue("benchFlowServices/monitors/mysql");
                    logger.info("mysqlMonitorEndpoint: " + mysqlMonitorEndpoint);
                } catch (XPathExpressionException ex) {
                    Thread t = Thread.currentThread();
                    t.getUncaughtExceptionHandler().uncaughtException(t, ex);
                    return;
                }

                String queryCall = "?query=SELECT+COUNT(*)+FROM+ACT_HI_PROCINST+WHERE+END_TIME_+IS+NULL&value=0&method=equal";
                //TODO: improve, the empty answer with equal seems no, the not empty seems yes in the current implementation
                String res = "";

                while(true) {
                    //TODO: we for sure want to have a better way to get the same.
                    //the point now is that it is not possible to throw an exception on the run method
                    try {
                        res = new BenchFlowServicesAsynchInteraction(mysqlMonitorEndpoint+queryCall).call();
                    } catch (Exception ex) {
                        Thread t = Thread.currentThread();
                        t.getUncaughtExceptionHandler().uncaughtException(t, ex);
                        return;
                    }

                    logger.info("Waiting for workload to complete, res: " + res);

                    //TODO: this is really custom given the current way the monitor works
                    if(res.toLowerCase().contains("matches 0")){
                        break;
                    } else {
                        //Pause for 10 seconds
                        logger.info("Waiting for workload to complete, waiting to restart");
                        //TODO: we for sure want to have a better way to get the same.
                        //the point now is that it is not possible to throw an exception on the run method
                        try {
                            Thread.sleep(10000);
                        } catch (InterruptedException ex) {
                            Thread t = Thread.currentThread();
                            t.getUncaughtExceptionHandler().uncaughtException(t, ex);
                            return;
                        }
                    }

                }

                done.countDown();
            }
        }).start();

        //Maybe the more precise point to stop the statistics collection is here?

        //TODO: we for sure want to have a better way to get the same.
        //the point now is that it is not possible to throw an exception on the run method
        try {
            //wait maximum 2 minutes (heuristically set)
            boolean processingCompleteWithin1Second = done.await(120000, TimeUnit.SECONDS);
        } catch (InterruptedException ex) {
            Thread t = Thread.currentThread();
            t.getUncaughtExceptionHandler().uncaughtException(t, ex);
            return;
        }

        logger.info("Post-run done");
    }

     /**
     * Do the request.
     * @throws IOException An I/O or network error occurred.
     */
    @BenchmarkOperation (
        name    = "activitiModel",
        max90th = 20000, // 20000 millisec
        timing  = Timing.AUTO
    )
    public void doActivitiModelRequest() throws IOException {

        if(isStarted()){
            String startURL = wfms.startProcessDefinition("activitiModel.bpmn");
        } else {
            String startURL = wfms.startProcessDefinition("mock.bpmn");
        }
        
    }

    /**
     * Do the request.
     * @throws IOException An I/O or network error occurred.
     */
    @BenchmarkOperation (
        name    = "camundaAdditionalApproval",
        max90th = 20000, // 20000 millisec
        timing  = Timing.AUTO
    )
    public void doCamundaAdditionalApprovalRequest() throws IOException {

        if(isStarted()){
            String startURL = wfms.startProcessDefinition("camundaAdditionalApproval.bpmn");
        } else {
            String startURL = wfms.startProcessDefinition("mock.bpmn");
        }
        
    }

    /**
     * Do the request.
     * @throws IOException An I/O or network error occurred.
     */
    @BenchmarkOperation (
        name    = "camundaCounterpartyOnboarding",
        max90th = 20000, // 20000 millisec
        timing  = Timing.AUTO
    )
    public void doCamundaCounterpartyOnboardingRequest() throws IOException {

        if(isStarted()){
            String startURL = wfms.startProcessDefinition("camundaCounterpartyOnboarding.bpmn");
        } else {
            String startURL = wfms.startProcessDefinition("mock.bpmn");
        }
        
    }

    /**
     * Do the request.
     * @throws IOException An I/O or network error occurred.
     */
    @BenchmarkOperation (
        name    = "camundaInvoiceCollaboration",
        max90th = 20000, // 20000 millisec
        timing  = Timing.AUTO
    )
    public void doCamundaInvoiceCollaborationRequest() throws IOException {

        if(isStarted()){
            String startURL = wfms.startProcessDefinition("camundaInvoiceCollaboration.bpmn");
        } else {
            String startURL = wfms.startProcessDefinition("mock.bpmn");
        }

    }

    /**
     * Do the request.
     * @throws IOException An I/O or network error occurred.
     */
    @BenchmarkOperation (
        name    = "camundaOop",
        max90th = 20000, // 20000 millisec
        timing  = Timing.AUTO
    )
    public void doCamundaOopRequest() throws IOException {

        if(isStarted()){
            String startURL = wfms.startProcessDefinition("camundaOop.bpmn");
        } else {
            String startURL = wfms.startProcessDefinition("mock.bpmn");
        }
        
    }

    /**
     * Things that can be abstracted away
     * Some of them must also be shared with the Benchmark definition, since they are basically duplicated code
     */
    private void initialize() {
        ctx = DriverContext.getContext();
        HttpTransport.setProvider("com.sun.faban.driver.transport.hc3.ApacheHC3Transport");
        http = HttpTransport.newInstance();
        logger = ctx.getLogger();
        //Needed to interact with the WfMS API
        JSONHeaders.put("Content-Type","application/json");
    }

    private String getContextProperty(String property){

        return ctx.getProperty(property);

    }

    private String getXPathValue(String xPathExpression) throws XPathExpressionException {

        return ctx.getXPathValue(xPathExpression);

    }

    private void setSUTEndpoint() throws XPathExpressionException {
        // SUTEndpoint = getXPathValue("sutConfig/SUTEndpoint");
        SUTEndpoint = getXPathValue("sutConfiguration/SUTEndpoint");
    }

    private void loadModelsInfo() {
        int numModel = Integer.parseInt(getContextProperty("model_num"));
        for (int i = 1; i <= numModel; i++) {
             String name = getContextProperty("model_" + i + "_name");
             String startID = getContextProperty("model_" + i + "_startID");
             // logger.warning("Model(" + i + ") name: " + name);
             // logger.warning("Model(" + i + ") startID: " + startID);
             modelsStartID.put(name, startID);
        }
    }

    /**
     *
     * TODO: improve, taken from the Faban codebase and customized
     *
     * @return true if this time span is in ramp up, steady state, ramp down, false otherwise.
     */
    boolean isStarted() {

        long steadyStateStartTime = ctx.getSteadyStateStartNanos();

        //If we don't have the steadyStateStartTime, it means it is not yet set,
        //then we are not during the run
        if( steadyStateStartTime!=0 ){

            long rampUpTime = ctx.getRampUp() * 1000000000l;
            long steadyStateTime = ctx.getSteadyState() * 1000000000l;
            long rampDownTime = ctx.getRampDown() * 1000000000l;
            
            long rampUpStartTime = steadyStateStartTime - rampUpTime;
            long steadyStateEndTime = steadyStateStartTime + steadyStateTime;
            long rampDownEndTime = steadyStateEndTime + rampDownTime;

            long currentTime = ctx.getNanoTime();
            
            logger.info("rampUpTime: " + rampUpTime);
            logger.info("steadyStateTime: " + steadyStateTime);
            logger.info("rampDownTime: " + rampDownTime);
            logger.info("rampUpStartTime: " + rampUpStartTime);
            logger.info("steadyStateEndTime: " + steadyStateEndTime);
            logger.info("rampDownEndTime: " + rampDownEndTime);
            logger.info("steadyStateStartTime: " + steadyStateStartTime);
            logger.info("currentTime: " + currentTime);

            return (rampUpStartTime <= currentTime) && (currentTime <= rampDownEndTime);
        }

        return false;
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
