/*
 * SPDX-License-Identifier: Apache-2.0
 */

'use strict';

const { FileSystemWallet, Gateway } = require('fabric-network');
const path = require('path');
const fs = require('fs')
const ccpPath = path.resolve(__dirname, '..', '..', 'first-network', 'connection-org1.json');
const jsonString = fs.readFileSync('/home/TIETRONIX.COM/jselvan/Desktop/test.json')

async function main() {
    try {

       

    
  
      
        const walletPath = path.join(process.cwd(), 'wallet');
        const wallet = new FileSystemWallet(walletPath);
        console.log(`Wallet path: ${walletPath}`);

       
        const userExists = await wallet.exists('user1');
        if (!userExists) {
            console.log('An identity for the user "admin" does not exist in the wallet');
            console.log('Run the registerUser.js application before retrying');
            return;
        }

 
        const gateway = new Gateway();
        await gateway.connect(ccpPath, { wallet, identity: 'user1', discovery: { enabled: true, asLocalhost: true } });

       
        const network = await gateway.getNetwork('mychannel');

       
        const contract = network.getContract('fabcar');

        //JSON transformation
        const temp = JSON.parse(jsonString)
        var myJSON = JSON.stringify(temp);
        
        await contract.submitTransaction('createComponent', myJSON);
        console.log('Transaction has been submitted'); 
       

        await gateway.disconnect();

    } catch (error) {
        console.error(`Failed to submit transaction: ${error}`);
        process.exit(1);
    }
}

main();
