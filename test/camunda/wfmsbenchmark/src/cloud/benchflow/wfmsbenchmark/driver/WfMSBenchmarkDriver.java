/*
 * Copyright (c) 2009-2010 Shanti Subramanyam, Akara Sucharitakul

 * Permission is hereby granted, free of charge, to any person
 * obtaining a copy of this software and associated documentation
 * files (the "Software"), to deal in the Software without
 * restriction, including without limitation the rights to use,
 * copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the
 * Software is furnished to do so, subject to the following
 * conditions:
 *
 * The above copyright notice and this permission notice shall be
 * included in all copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
 * EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES
 * OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
 * NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT
 * HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY,
 * WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
 * FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR
 * OTHER DEALINGS IN THE SOFTWARE.
 */
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

/**
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
    operations = {"parallelSplitPattern", "exclusiveChoicePattern", "implicitTerminationPattern"},
    mix = { @Row({ 50, 50, 50 }),
            @Row({ 50, 50, 50 }),
            @Row({ 50, 50, 50 }) },
    deviation = 5
)
// @MatrixMix (
//     operations = {"parallelSplitPattern", "exclusiveChoicePattern"},
//     mix = { @Row({ 50, 50 }),
//             @Row({ 50, 50 }) },
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


    /**
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
    /**
     * Start monitors and collectors
     * Only thread 0 does is and can share data with onceAfter
     */
    @OnceBefore
    public void onceBefore() throws XPathExpressionException, IOException {
        // Stats:
        // - curl -v -X GET http://10.40.1.128:9302/start
        String statsStart = getXPathValue("services/collectors/statsStart");
        http.fetchURL(statsStart);
    }

    /**
     * Ask monitors to wait for completion
     * Only thread 0 does is and can share data with onceBefore
     */
    @OnceAfter
    public void onceAfter() {
        //Currenlty we just wait
        try {
            Thread.sleep(60000);
        } catch (InterruptedException e) {
        }
        logger.info("Tested post-run (sleep 60) done");
    }

    /**
     * Do the request.
     * @throws IOException An I/O or network error occurred.
     */
    @BenchmarkOperation (
        name    = "parallelSplitPattern",
        max90th = 20000, // 20000 millisec
        timing  = Timing.AUTO
    )
    public void doParallelSplitPatternRequest() throws IOException {

        String startURL = wfms.startProcessDefinition("parallelSplitPattern.bpmn");

    }

    /**
     * Do the request.
     * @throws IOException An I/O or network error occurred.
     */
    @BenchmarkOperation (
        name    = "exclusiveChoicePattern",
        max90th = 20000, // 20000 millisec
        timing  = Timing.AUTO
    )
    public void doExclusiveChoicePatternRequest() throws IOException {

        String startURL = wfms.startProcessDefinition("exclusiveChoicePattern.bpmn");
        
    }

    /**
     * Do the request.
     * @throws IOException An I/O or network error occurred.
     */
    @BenchmarkOperation (
        name    = "implicitTerminationPattern",
        max90th = 20000, // 20000 millisec
        timing  = Timing.AUTO
    )
    public void doImplicitTerminationPatternRequest() throws IOException {

        String startURL = wfms.startProcessDefinition("implicitTerminationPattern.bpmn");
        
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

        SUTEndpoint = getXPathValue("sutConfig/SUTEndpoint");

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
}
