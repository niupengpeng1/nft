import { network } from "hardhat";

const { ethers } = await network.connect({
    network: "hardhatOp",
    chainType: "op",
});

async function main() {

    const [deployer] = await ethers.getSigners();
    const balance = await ethers.provider.getBalance(deployer.address);
    console.log("Deploying contracts with account:", deployer.address);
    console.log("Account balance:", ethers.formatEther(balance), "ETH");

    if (balance === 0n) {
        throw new Error("Account has no balance. Please fund the account first.");
    }

    const counter = await ethers.deployContract("Counter");

    console.log("\nDeploying Counter contract...");

    await counter.waitForDeployment();

    const code = await ethers.provider.getCode(await counter.getAddress());
    if (code === "0x") {
        throw new Error("Contract deployment failed. Please check the transaction receipt.");
    }


    console.log("Counter contract deployed to:", await counter.getAddress());

    await counter.inc();

    await counter.incBy(123);

    const nowValue = await counter.get();
    console.log("Current value:", nowValue.toString());

    console.log("Deployment complete.");
}


// 执行主函数
main()
    .then(() => process.exit(0))
    .catch((error) => {
        console.error(error);
        process.exit(1);
    });
