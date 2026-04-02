import { network } from "hardhat";

const { ethers } = await network.connect({
    network: "hardhatOp",
    chainType: "op",
});


async function main() {
    const [deployer] = await ethers.getSigners();

    console.log("Deploying contracts with account:", deployer.address);
    console.log("Account balance:", ethers.formatEther(await ethers.provider.getBalance(deployer.address)), "ETH");
    const hellow = await ethers.deployContract("HelloWorld",["你好，世界！"]);

    await hellow.waitForDeployment();

    console.log("Contract deployed to:", await hellow.getAddress());

    console.log("更新前Greeting:", await hellow.getMessage());

    await hellow.updateMessage("Hello, Hardhat!");
    console.log("更新后Greeting:", await hellow.getMessage());


     // ==========================================
  // 👇 下面这一段就是查看 Gas 消耗的代码
  // ==========================================
   const deploymentTx = hellow.deploymentTransaction();
    if (deploymentTx) {
        // 等待交易收据
        const deploymentReceipt = await deploymentTx.wait();
        
        if (deploymentReceipt) {
            const gasUsed = deploymentReceipt.gasUsed;
            const gasPrice =deploymentReceipt.gasPrice;

            console.log("✅ 部署成功");
            console.log("Gas 消耗:", gasUsed.toString());
            console.log("Gas 价格:", gasPrice?.toString() || 'N/A');

            if (gasPrice) {
                const totalWei = gasUsed * gasPrice;
                console.log("总花费(wei):", totalWei.toString());
                console.log("总花费(ether):", ethers.formatEther(totalWei));
            }
        } else {
            console.log("⚠️ 部署交易收据不可用");
        }
    } else {
        console.log("⚠️ 部署交易不可用");
    }
}


main()
    .then(() => process.exit(0))
    .catch((error) => {
        console.error(error);
        process.exit(1);
    });