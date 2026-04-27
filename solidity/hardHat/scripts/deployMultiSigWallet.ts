import  hre  from "hardhat";

const { ethers } = await hre.network.getOrCreate();

async function main() {

    const [deployer,admin] = await ethers.getSigners();
    const balance = await ethers.provider.getBalance(deployer.address);
    console.log("Deploying contracts with account:", deployer.address);
    console.log("Account balance:", ethers.formatEther(balance), "ETH");

    if (balance === 0n) {
        throw new Error("Account has no balance. Please fund the account first.");
    }

    const MultiSigWallet = await ethers.getContractFactory("MultiSigWallet");
    const multiSigWallet = await MultiSigWallet.deploy([deployer,admin],2);

    console.log("\nDeploying MultiSigWallet contract...");

    await multiSigWallet.waitForDeployment();

    const code = await ethers.provider.getCode(await multiSigWallet.getAddress());
    if (code === "0x") {
        throw new Error("multiSigWallet deployment failed. Please check the transaction receipt.");
    }


    console.log("MultiSigWallet contract deployed to:", await multiSigWallet.getAddress());

    
    

    console.log("Deployment complete.");
}

// 执行主函数
main()
    .then(() => process.exit(0))
    .catch((error) => {
        console.error(error);
        process.exit(1);
    });