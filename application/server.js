// ExpressJS Setup
const express = require("express");
const app = express();
var bodyParser = require("body-parser");

// Hyperledger Bridge
const { Wallets, Gateway } = require("fabric-network");
const fs = require("fs");
const path = require("path");
const ccpPath = path.resolve(__dirname, "ccp", "connection-org1.json");
let ccp = JSON.parse(fs.readFileSync(ccpPath, "utf8"));

// Constants
const PORT = 8080;
const HOST = "0.0.0.0";

// use static file
app.use(express.static(path.join(__dirname, "views")));

// configure app to use body-parser
app.use(bodyParser.json());
app.use(bodyParser.urlencoded({ extended: false }));

// main page routing
app.get("/", (req, res) => {
    res.sendFile(__dirname + "/index.html");
});

async function cc_call(fn_name, args) {
    const walletPath = path.join(process.cwd(), "wallet");
    console.log(`Wallet path: ${walletPath}`);
    const wallet = await Wallets.newFileSystemWallet(walletPath);

    console.log(`Wallet path: ${walletPath}`);

    const userExists = await wallet.get("appUser");
    if (!userExists) {
        console.log(
            'An identity for the user "appUser" does not exist in the wallet'
        );
        console.log("Run the registerUser.js application before retrying");
        return;
    }

    const gateway = new Gateway();
    await gateway.connect(ccp, {
        wallet,
        identity: "appUser",
        discovery: { enabled: true, asLocalhost: true },
    });

    const network = await gateway.getNetwork("mychannel");
    const contract = network.getContract("tablet");

    var result;

    if (fn_name == "addUser") {
        result = await contract.submitTransaction("addUser", args);
    } 
    else if (fn_name == "addDate") {
        e = args[0];
        p = args[1];
        s = args[2];
        result = await contract.submitTransaction("addDate", e, p, s);
    } 
    else if (fn_name == "returnTablet") {
        e = args[0];
        p = args[1];
        s = args[2];
        result = await contract.submitTransaction("returnTablet", e, p, s);
    } 
    else if (fn_name == "readDate")
        result = await contract.evaluateTransaction("readDate", args);
    else result = "not supported function";

    return result;
}

// create mate
app.post("/mate", async (req, res) => {
    const email = req.body.email;
    console.log("add user email: " + email);

    //result = await cc_call("addUser", email);
    result = await cc_call("addUser", email);

    const myobj = { result: "success" };
    res.status(200).json(myobj);
});


// add date
app.post("/record", async (req, res) => {
    const email = req.body.email;
    const tablet = req.body.tablet;
    const record = req.body.record;
    console.log("add email: " + email);
    console.log("add tablet: " + tablet);
    console.log("add record: " + record);

    var args = [email, tablet, record];
    result = await cc_call("addDate", args);

    const myobj = { result: "success" };
    res.status(200).json(myobj);
});

// return tablet
app.post("/return", async (req, res) => {
    const email = req.body.email;
    const tablet = req.body.tablet;
    const record = req.body.record;
    console.log("return email: " + email);
    console.log("return tablet: " + tablet);
    console.log("return record: " + record);

    var args = [email, tablet, record];
    result = await cc_call("returnTablet", args);

    const myobj = { result: "success" };
    res.status(200).json(myobj);
});

// find mate
app.post("/mate/:email", async (req, res) => {
    const email = req.body.email;
    console.log("email: " + req.body.email);

    result = await cc_call("readDate", email);
    console.log("result: " + result.toString());
    const myobj = JSON.parse(result, "utf8");

    res.status(200).json(myobj);
});

// server start
app.listen(PORT, HOST);
console.log(`Running on http://${HOST}:${PORT}`);
