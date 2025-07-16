<<<<<<< HEAD
import { ethers } from "ethers";

// Operator Clusters Registry Contract
const ObolOperatorClustersRegistry = "0xaE7B191A31f627b4eB1d4DaC64eaB9976995b433";

// Obol Operator Clusters Registry Contract ABI
const ObolOperatorClustersRegistryProxyABI = [{"constant":true,"inputs":[],"name":"proxyType","outputs":[{"name":"proxyTypeId","type":"uint256"}],"payable":false,"stateMutability":"pure","type":"function"},{"constant":true,"inputs":[],"name":"isDepositable","outputs":[{"name":"","type":"bool"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[],"name":"implementation","outputs":[{"name":"","type":"address"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[],"name":"appId","outputs":[{"name":"","type":"bytes32"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[],"name":"kernel","outputs":[{"name":"","type":"address"}],"payable":false,"stateMutability":"view","type":"function"},{"inputs":[{"name":"_kernel","type":"address"},{"name":"_appId","type":"bytes32"},{"name":"_initializePayload","type":"bytes"}],"payable":false,"stateMutability":"nonpayable","type":"constructor"},{"payable":true,"stateMutability":"payable","type":"fallback"},{"anonymous":false,"inputs":[{"indexed":false,"name":"sender","type":"address"},{"indexed":false,"name":"value","type":"uint256"}],"name":"ProxyDeposit","type":"event"}]
const ObolOperatorClustersRegistryLogicABI = [{"constant":true,"inputs":[],"name":"getActiveNodeOperatorsCount","outputs":[{"name":"","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[{"name":"_nodeOperatorId","type":"uint256"},{"name":"_fullInfo","type":"bool"}],"name":"getNodeOperator","outputs":[{"name":"active","type":"bool"},{"name":"name","type":"string"},{"name":"rewardAddress","type":"address"},{"name":"totalVettedValidators","type":"uint64"},{"name":"totalExitedValidators","type":"uint64"},{"name":"totalAddedValidators","type":"uint64"},{"name":"totalDepositedValidators","type":"uint64"}],"payable":false,"stateMutability":"view","type":"function"}];
const ObolOperatorABI = [{"inputs":[{"internalType":"address","name":"_feeRecipient","type":"address"},{"internalType":"uint256","name":"_feeShare","type":"uint256"},{"internalType":"contract ERC20","name":"_stETH","type":"address"},{"internalType":"contract ERC20","name":"_wstETH","type":"address"}],"stateMutability":"nonpayable","type":"constructor"},{"inputs":[],"name":"Invalid_Address","type":"error"},{"inputs":[],"name":"Invalid_FeeRecipient","type":"error"},{"inputs":[{"internalType":"uint256","name":"fee","type":"uint256"}],"name":"Invalid_FeeShare","type":"error"},{"inputs":[],"name":"distribute","outputs":[{"internalType":"uint256","name":"amount","type":"uint256"}],"stateMutability":"nonpayable","type":"function"},{"inputs":[],"name":"feeRecipient","outputs":[{"internalType":"address","name":"","type":"address"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"feeShare","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"address","name":"token","type":"address"}],"name":"rescueFunds","outputs":[{"internalType":"uint256","name":"balance","type":"uint256"}],"stateMutability":"nonpayable","type":"function"},{"inputs":[],"name":"splitWallet","outputs":[{"internalType":"address","name":"","type":"address"}],"stateMutability":"pure","type":"function"},{"inputs":[],"name":"stETH","outputs":[{"internalType":"contract ERC20","name":"","type":"address"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"wstETH","outputs":[{"internalType":"contract ERC20","name":"","type":"address"}],"stateMutability":"view","type":"function"}];
const fullABI = [...ObolOperatorClustersRegistryProxyABI, ...ObolOperatorClustersRegistryLogicABI];
const SplitMainInterface = new ethers.Interface([{"inputs":[],"stateMutability":"nonpayable","type":"constructor"},{"inputs":[],"name":"Create2Error","type":"error"},{"inputs":[],"name":"CreateError","type":"error"},{"inputs":[{"internalType":"address","name":"newController","type":"address"}],"name":"InvalidNewController","type":"error"},{"inputs":[{"internalType":"uint256","name":"accountsLength","type":"uint256"},{"internalType":"uint256","name":"allocationsLength","type":"uint256"}],"name":"InvalidSplit__AccountsAndAllocationsMismatch","type":"error"},{"inputs":[{"internalType":"uint256","name":"index","type":"uint256"}],"name":"InvalidSplit__AccountsOutOfOrder","type":"error"},{"inputs":[{"internalType":"uint256","name":"index","type":"uint256"}],"name":"InvalidSplit__AllocationMustBePositive","type":"error"},{"inputs":[{"internalType":"uint32","name":"allocationsSum","type":"uint32"}],"name":"InvalidSplit__InvalidAllocationsSum","type":"error"},{"inputs":[{"internalType":"uint32","name":"distributorFee","type":"uint32"}],"name":"InvalidSplit__InvalidDistributorFee","type":"error"},{"inputs":[{"internalType":"bytes32","name":"hash","type":"bytes32"}],"name":"InvalidSplit__InvalidHash","type":"error"},{"inputs":[{"internalType":"uint256","name":"accountsLength","type":"uint256"}],"name":"InvalidSplit__TooFewAccounts","type":"error"},{"inputs":[{"internalType":"address","name":"sender","type":"address"}],"name":"Unauthorized","type":"error"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"address","name":"split","type":"address"}],"name":"CancelControlTransfer","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"address","name":"split","type":"address"},{"indexed":true,"internalType":"address","name":"previousController","type":"address"},{"indexed":true,"internalType":"address","name":"newController","type":"address"}],"name":"ControlTransfer","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"address","name":"split","type":"address"}],"name":"CreateSplit","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"address","name":"split","type":"address"},{"indexed":true,"internalType":"contract ERC20","name":"token","type":"address"},{"indexed":false,"internalType":"uint256","name":"amount","type":"uint256"},{"indexed":true,"internalType":"address","name":"distributorAddress","type":"address"}],"name":"DistributeERC20","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"address","name":"split","type":"address"},{"indexed":false,"internalType":"uint256","name":"amount","type":"uint256"},{"indexed":true,"internalType":"address","name":"distributorAddress","type":"address"}],"name":"DistributeETH","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"address","name":"split","type":"address"},{"indexed":true,"internalType":"address","name":"newPotentialController","type":"address"}],"name":"InitiateControlTransfer","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"address","name":"split","type":"address"}],"name":"UpdateSplit","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"address","name":"account","type":"address"},{"indexed":false,"internalType":"uint256","name":"ethAmount","type":"uint256"},{"indexed":false,"internalType":"contract ERC20[]","name":"tokens","type":"address[]"},{"indexed":false,"internalType":"uint256[]","name":"tokenAmounts","type":"uint256[]"}],"name":"Withdrawal","type":"event"},{"inputs":[],"name":"PERCENTAGE_SCALE","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"address","name":"split","type":"address"}],"name":"acceptControl","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"address","name":"split","type":"address"}],"name":"cancelControlTransfer","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"address[]","name":"accounts","type":"address[]"},{"internalType":"uint32[]","name":"percentAllocations","type":"uint32[]"},{"internalType":"uint32","name":"distributorFee","type":"uint32"},{"internalType":"address","name":"controller","type":"address"}],"name":"createSplit","outputs":[{"internalType":"address","name":"split","type":"address"}],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"address","name":"split","type":"address"},{"internalType":"contract ERC20","name":"token","type":"address"},{"internalType":"address[]","name":"accounts","type":"address[]"},{"internalType":"uint32[]","name":"percentAllocations","type":"uint32[]"},{"internalType":"uint32","name":"distributorFee","type":"uint32"},{"internalType":"address","name":"distributorAddress","type":"address"}],"name":"distributeERC20","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"address","name":"split","type":"address"},{"internalType":"address[]","name":"accounts","type":"address[]"},{"internalType":"uint32[]","name":"percentAllocations","type":"uint32[]"},{"internalType":"uint32","name":"distributorFee","type":"uint32"},{"internalType":"address","name":"distributorAddress","type":"address"}],"name":"distributeETH","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"address","name":"split","type":"address"}],"name":"getController","outputs":[{"internalType":"address","name":"","type":"address"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"address","name":"account","type":"address"},{"internalType":"contract ERC20","name":"token","type":"address"}],"name":"getERC20Balance","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"address","name":"account","type":"address"}],"name":"getETHBalance","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"address","name":"split","type":"address"}],"name":"getHash","outputs":[{"internalType":"bytes32","name":"","type":"bytes32"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"address","name":"split","type":"address"}],"name":"getNewPotentialController","outputs":[{"internalType":"address","name":"","type":"address"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"address","name":"split","type":"address"}],"name":"makeSplitImmutable","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"address[]","name":"accounts","type":"address[]"},{"internalType":"uint32[]","name":"percentAllocations","type":"uint32[]"},{"internalType":"uint32","name":"distributorFee","type":"uint32"}],"name":"predictImmutableSplitAddress","outputs":[{"internalType":"address","name":"split","type":"address"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"address","name":"split","type":"address"},{"internalType":"address","name":"newController","type":"address"}],"name":"transferControl","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"address","name":"split","type":"address"},{"internalType":"contract ERC20","name":"token","type":"address"},{"internalType":"address[]","name":"accounts","type":"address[]"},{"internalType":"uint32[]","name":"percentAllocations","type":"uint32[]"},{"internalType":"uint32","name":"distributorFee","type":"uint32"},{"internalType":"address","name":"distributorAddress","type":"address"}],"name":"updateAndDistributeERC20","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"address","name":"split","type":"address"},{"internalType":"address[]","name":"accounts","type":"address[]"},{"internalType":"uint32[]","name":"percentAllocations","type":"uint32[]"},{"internalType":"uint32","name":"distributorFee","type":"uint32"},{"internalType":"address","name":"distributorAddress","type":"address"}],"name":"updateAndDistributeETH","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"address","name":"split","type":"address"},{"internalType":"address[]","name":"accounts","type":"address[]"},{"internalType":"uint32[]","name":"percentAllocations","type":"uint32[]"},{"internalType":"uint32","name":"distributorFee","type":"uint32"}],"name":"updateSplit","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[],"name":"walletImplementation","outputs":[{"internalType":"address","name":"","type":"address"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"address","name":"account","type":"address"},{"internalType":"uint256","name":"withdrawETH","type":"uint256"},{"internalType":"contract ERC20[]","name":"tokens","type":"address[]"}],"name":"withdraw","outputs":[],"stateMutability":"nonpayable","type":"function"},{"stateMutability":"payable","type":"receive"}]);



const EtherscanAPIKey = "RKFRN7F2EGCMFWPMVTX5RE4ZIQ7E1RHJXI";
const RPC_URL = "https://convincing-divine-slug.quiknode.pro/9e835fa839e702f4140f3a80b65672cfb3d5950b/";
const provider = new ethers.JsonRpcProvider(RPC_URL);
const etherscanProvider = new ethers.EtherscanProvider('homestead',EtherscanAPIKey);
const ObolOperatorClustersRegistryContract = new ethers.Contract(ObolOperatorClustersRegistry, fullABI, provider);


// 請求節流配置
const BATCH_SIZE = 8; // 每批最多10個請求
const BATCH_DELAY = 1000; // 批次間延遲1秒

// 延遲函數
const delay = (ms) => new Promise(resolve => setTimeout(resolve, ms));

// 分批處理請求
const processBatches = async (items, batchProcessor) => {
  const results = [];
  
  for (let i = 0; i < items.length; i += BATCH_SIZE) {
    const batch = items.slice(i, i + BATCH_SIZE);
    console.log(`Processing batch ${Math.floor(i / BATCH_SIZE) + 1}/${Math.ceil(items.length / BATCH_SIZE)}: ${batch.length} items`);
    
    try {
      const batchResults = await batchProcessor(batch);
      results.push(...batchResults);
    } catch (error) {
      console.error(`Batch ${Math.floor(i / BATCH_SIZE) + 1} failed:`, error);
      // 添加空結果以保持索引一致
      results.push(...new Array(batch.length).fill(null));
    }
    
    // 如果不是最後一批，等待一段時間
    if (i + BATCH_SIZE < items.length) {
      console.log(`Waiting ${BATCH_DELAY}ms before next batch...`);
      await delay(BATCH_DELAY);
    }
  }
  
  return results;
};

export const ether_obol = {
  getObolOperatorClustersRegistry: async () => {
    try {
      console.log('Fetching active node operators count...');
      const activeNodeOperators = await ObolOperatorClustersRegistryContract.getActiveNodeOperatorsCount()
      const operatorCount = Number(activeNodeOperators) // 轉換 BigNumber 為普通數字
      console.log('Active node operators count:', operatorCount)
      
      if (operatorCount === 0) {
        console.log('No active node operators found')
        return []
      }
      
      // 創建操作者索引數組
      const operatorIndices = Array.from({ length: operatorCount }, (_, i) => i);
      
      // 批次處理器函數
      const batchProcessor = async (indices) => {
        const promises = indices.map(index => 
          ObolOperatorClustersRegistryContract.getNodeOperator(index, true).catch(error => {
            console.error(`Error fetching operator ${index}:`, error);
            return null;
          })
        );
        
        return Promise.all(promises);
      };
      
      console.log(`Starting to fetch ${operatorCount} operators in batches of ${BATCH_SIZE}...`);
      const nodeOperators = await processBatches(operatorIndices, batchProcessor);
      
      // 過濾掉失敗的調用
      const validOperators = nodeOperators.filter(operator => operator !== null);
      console.log(`Successfully fetched ${validOperators.length}/${operatorCount} operators`);
      
      return validOperators;
    } catch (error) {
      console.error('Error in getObolOperatorClustersRegistry:', error);
      throw error;
    }
  },
  getObolOperatorSplitWallets: async (operatorAddress) => {
    const operatorContract = new ethers.Contract(operatorAddress, ObolOperatorABI, provider);
    const splitWallet = await operatorContract.splitWallet();
    return splitWallet;
  },
  getObolOperatorRewardshare: async (splitWallet) => {
    const CREATE_SELECTOR = "0x3d605d80";
    const url = `https://api.etherscan.io/api?module=account&action=txlistinternal&address=${splitWallet}&startblock=0&endblock=99999999&sort=desc&apikey=${EtherscanAPIKey}`;
    console.log('Fetching transactions for split wallet:', splitWallet);
    const transactions = await fetch(url).then(res => res.json()).then(data => data.result);
    console.log('Transactions:', transactions);
    const createTx = transactions.find(tx => tx.type == "create2")
    console.log('Create tx:', createTx);
    const createOriginTx = await etherscanProvider.getTransaction(createTx.hash);
    console.log('Create origin tx:', createOriginTx.data);
    const decodedData = SplitMainInterface.decodeFunctionData('createSplit', createOriginTx.data);
    console.log('Decoded data:', decodedData);
    if (decodedData) {
      const rewardAddress = decodedData[0];
      const rewardShare = decodedData[1];
      console.log('Reward address:', rewardAddress);
      console.log('Reward share:', rewardShare);
      return {rewardAddress,rewardShare};
    } else {
      throw new Error('Failed to decode createSplit function data');
    }
  }
=======
import { ethers } from "ethers";

// Operator Clusters Registry Contract
const ObolOperatorClustersRegistry = "0xaE7B191A31f627b4eB1d4DaC64eaB9976995b433";

// Obol Operator Clusters Registry Contract ABI
const ObolOperatorClustersRegistryProxyABI = [{"constant":true,"inputs":[],"name":"proxyType","outputs":[{"name":"proxyTypeId","type":"uint256"}],"payable":false,"stateMutability":"pure","type":"function"},{"constant":true,"inputs":[],"name":"isDepositable","outputs":[{"name":"","type":"bool"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[],"name":"implementation","outputs":[{"name":"","type":"address"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[],"name":"appId","outputs":[{"name":"","type":"bytes32"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[],"name":"kernel","outputs":[{"name":"","type":"address"}],"payable":false,"stateMutability":"view","type":"function"},{"inputs":[{"name":"_kernel","type":"address"},{"name":"_appId","type":"bytes32"},{"name":"_initializePayload","type":"bytes"}],"payable":false,"stateMutability":"nonpayable","type":"constructor"},{"payable":true,"stateMutability":"payable","type":"fallback"},{"anonymous":false,"inputs":[{"indexed":false,"name":"sender","type":"address"},{"indexed":false,"name":"value","type":"uint256"}],"name":"ProxyDeposit","type":"event"}]
const ObolOperatorClustersRegistryLogicABI = [{"constant":true,"inputs":[],"name":"getActiveNodeOperatorsCount","outputs":[{"name":"","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[{"name":"_nodeOperatorId","type":"uint256"},{"name":"_fullInfo","type":"bool"}],"name":"getNodeOperator","outputs":[{"name":"active","type":"bool"},{"name":"name","type":"string"},{"name":"rewardAddress","type":"address"},{"name":"totalVettedValidators","type":"uint64"},{"name":"totalExitedValidators","type":"uint64"},{"name":"totalAddedValidators","type":"uint64"},{"name":"totalDepositedValidators","type":"uint64"}],"payable":false,"stateMutability":"view","type":"function"}];
const ObolOperatorABI = [{"inputs":[{"internalType":"address","name":"_feeRecipient","type":"address"},{"internalType":"uint256","name":"_feeShare","type":"uint256"},{"internalType":"contract ERC20","name":"_stETH","type":"address"},{"internalType":"contract ERC20","name":"_wstETH","type":"address"}],"stateMutability":"nonpayable","type":"constructor"},{"inputs":[],"name":"Invalid_Address","type":"error"},{"inputs":[],"name":"Invalid_FeeRecipient","type":"error"},{"inputs":[{"internalType":"uint256","name":"fee","type":"uint256"}],"name":"Invalid_FeeShare","type":"error"},{"inputs":[],"name":"distribute","outputs":[{"internalType":"uint256","name":"amount","type":"uint256"}],"stateMutability":"nonpayable","type":"function"},{"inputs":[],"name":"feeRecipient","outputs":[{"internalType":"address","name":"","type":"address"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"feeShare","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"address","name":"token","type":"address"}],"name":"rescueFunds","outputs":[{"internalType":"uint256","name":"balance","type":"uint256"}],"stateMutability":"nonpayable","type":"function"},{"inputs":[],"name":"splitWallet","outputs":[{"internalType":"address","name":"","type":"address"}],"stateMutability":"pure","type":"function"},{"inputs":[],"name":"stETH","outputs":[{"internalType":"contract ERC20","name":"","type":"address"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"wstETH","outputs":[{"internalType":"contract ERC20","name":"","type":"address"}],"stateMutability":"view","type":"function"}];
const fullABI = [...ObolOperatorClustersRegistryProxyABI, ...ObolOperatorClustersRegistryLogicABI];
const SplitMainInterface = new ethers.Interface([{"inputs":[],"stateMutability":"nonpayable","type":"constructor"},{"inputs":[],"name":"Create2Error","type":"error"},{"inputs":[],"name":"CreateError","type":"error"},{"inputs":[{"internalType":"address","name":"newController","type":"address"}],"name":"InvalidNewController","type":"error"},{"inputs":[{"internalType":"uint256","name":"accountsLength","type":"uint256"},{"internalType":"uint256","name":"allocationsLength","type":"uint256"}],"name":"InvalidSplit__AccountsAndAllocationsMismatch","type":"error"},{"inputs":[{"internalType":"uint256","name":"index","type":"uint256"}],"name":"InvalidSplit__AccountsOutOfOrder","type":"error"},{"inputs":[{"internalType":"uint256","name":"index","type":"uint256"}],"name":"InvalidSplit__AllocationMustBePositive","type":"error"},{"inputs":[{"internalType":"uint32","name":"allocationsSum","type":"uint32"}],"name":"InvalidSplit__InvalidAllocationsSum","type":"error"},{"inputs":[{"internalType":"uint32","name":"distributorFee","type":"uint32"}],"name":"InvalidSplit__InvalidDistributorFee","type":"error"},{"inputs":[{"internalType":"bytes32","name":"hash","type":"bytes32"}],"name":"InvalidSplit__InvalidHash","type":"error"},{"inputs":[{"internalType":"uint256","name":"accountsLength","type":"uint256"}],"name":"InvalidSplit__TooFewAccounts","type":"error"},{"inputs":[{"internalType":"address","name":"sender","type":"address"}],"name":"Unauthorized","type":"error"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"address","name":"split","type":"address"}],"name":"CancelControlTransfer","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"address","name":"split","type":"address"},{"indexed":true,"internalType":"address","name":"previousController","type":"address"},{"indexed":true,"internalType":"address","name":"newController","type":"address"}],"name":"ControlTransfer","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"address","name":"split","type":"address"}],"name":"CreateSplit","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"address","name":"split","type":"address"},{"indexed":true,"internalType":"contract ERC20","name":"token","type":"address"},{"indexed":false,"internalType":"uint256","name":"amount","type":"uint256"},{"indexed":true,"internalType":"address","name":"distributorAddress","type":"address"}],"name":"DistributeERC20","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"address","name":"split","type":"address"},{"indexed":false,"internalType":"uint256","name":"amount","type":"uint256"},{"indexed":true,"internalType":"address","name":"distributorAddress","type":"address"}],"name":"DistributeETH","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"address","name":"split","type":"address"},{"indexed":true,"internalType":"address","name":"newPotentialController","type":"address"}],"name":"InitiateControlTransfer","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"address","name":"split","type":"address"}],"name":"UpdateSplit","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"address","name":"account","type":"address"},{"indexed":false,"internalType":"uint256","name":"ethAmount","type":"uint256"},{"indexed":false,"internalType":"contract ERC20[]","name":"tokens","type":"address[]"},{"indexed":false,"internalType":"uint256[]","name":"tokenAmounts","type":"uint256[]"}],"name":"Withdrawal","type":"event"},{"inputs":[],"name":"PERCENTAGE_SCALE","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"address","name":"split","type":"address"}],"name":"acceptControl","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"address","name":"split","type":"address"}],"name":"cancelControlTransfer","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"address[]","name":"accounts","type":"address[]"},{"internalType":"uint32[]","name":"percentAllocations","type":"uint32[]"},{"internalType":"uint32","name":"distributorFee","type":"uint32"},{"internalType":"address","name":"controller","type":"address"}],"name":"createSplit","outputs":[{"internalType":"address","name":"split","type":"address"}],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"address","name":"split","type":"address"},{"internalType":"contract ERC20","name":"token","type":"address"},{"internalType":"address[]","name":"accounts","type":"address[]"},{"internalType":"uint32[]","name":"percentAllocations","type":"uint32[]"},{"internalType":"uint32","name":"distributorFee","type":"uint32"},{"internalType":"address","name":"distributorAddress","type":"address"}],"name":"distributeERC20","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"address","name":"split","type":"address"},{"internalType":"address[]","name":"accounts","type":"address[]"},{"internalType":"uint32[]","name":"percentAllocations","type":"uint32[]"},{"internalType":"uint32","name":"distributorFee","type":"uint32"},{"internalType":"address","name":"distributorAddress","type":"address"}],"name":"distributeETH","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"address","name":"split","type":"address"}],"name":"getController","outputs":[{"internalType":"address","name":"","type":"address"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"address","name":"account","type":"address"},{"internalType":"contract ERC20","name":"token","type":"address"}],"name":"getERC20Balance","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"address","name":"account","type":"address"}],"name":"getETHBalance","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"address","name":"split","type":"address"}],"name":"getHash","outputs":[{"internalType":"bytes32","name":"","type":"bytes32"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"address","name":"split","type":"address"}],"name":"getNewPotentialController","outputs":[{"internalType":"address","name":"","type":"address"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"address","name":"split","type":"address"}],"name":"makeSplitImmutable","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"address[]","name":"accounts","type":"address[]"},{"internalType":"uint32[]","name":"percentAllocations","type":"uint32[]"},{"internalType":"uint32","name":"distributorFee","type":"uint32"}],"name":"predictImmutableSplitAddress","outputs":[{"internalType":"address","name":"split","type":"address"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"address","name":"split","type":"address"},{"internalType":"address","name":"newController","type":"address"}],"name":"transferControl","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"address","name":"split","type":"address"},{"internalType":"contract ERC20","name":"token","type":"address"},{"internalType":"address[]","name":"accounts","type":"address[]"},{"internalType":"uint32[]","name":"percentAllocations","type":"uint32[]"},{"internalType":"uint32","name":"distributorFee","type":"uint32"},{"internalType":"address","name":"distributorAddress","type":"address"}],"name":"updateAndDistributeERC20","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"address","name":"split","type":"address"},{"internalType":"address[]","name":"accounts","type":"address[]"},{"internalType":"uint32[]","name":"percentAllocations","type":"uint32[]"},{"internalType":"uint32","name":"distributorFee","type":"uint32"},{"internalType":"address","name":"distributorAddress","type":"address"}],"name":"updateAndDistributeETH","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"address","name":"split","type":"address"},{"internalType":"address[]","name":"accounts","type":"address[]"},{"internalType":"uint32[]","name":"percentAllocations","type":"uint32[]"},{"internalType":"uint32","name":"distributorFee","type":"uint32"}],"name":"updateSplit","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[],"name":"walletImplementation","outputs":[{"internalType":"address","name":"","type":"address"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"address","name":"account","type":"address"},{"internalType":"uint256","name":"withdrawETH","type":"uint256"},{"internalType":"contract ERC20[]","name":"tokens","type":"address[]"}],"name":"withdraw","outputs":[],"stateMutability":"nonpayable","type":"function"},{"stateMutability":"payable","type":"receive"}]);



const EtherscanAPIKey = "RKFRN7F2EGCMFWPMVTX5RE4ZIQ7E1RHJXI";
const RPC_URL = "https://convincing-divine-slug.quiknode.pro/9e835fa839e702f4140f3a80b65672cfb3d5950b/";
const provider = new ethers.JsonRpcProvider(RPC_URL);
const etherscanProvider = new ethers.EtherscanProvider('homestead',EtherscanAPIKey);
const ObolOperatorClustersRegistryContract = new ethers.Contract(ObolOperatorClustersRegistry, fullABI, provider);
const SECONDS_PER_BLOCK = 12;              // 主網平均
const BLOCKS_PER_DAY    = 86_400 / SECONDS_PER_BLOCK;


// 請求節流配置
const BATCH_SIZE = 8; // 每批最多8個請求
const BATCH_DELAY = 1000; // 批次間延遲1秒

// 圖表數據專用的併發配置（適應 QuickNode 15 req/sec 限制）
const CHART_BATCH_SIZE = 6; // 圖表數據每批最多6個區塊（12個請求：6個 getNodeOperator + 6個 getBlock）
const CHART_BATCH_DELAY = 1000; // 圖表數據批次間延遲1秒，確保不超過 14 req/sec

// 延遲函數
const delay = (ms) => new Promise(resolve => setTimeout(resolve, ms));

// 分批處理請求
const processBatches = async (items, batchProcessor) => {
  const results = [];
  
  for (let i = 0; i < items.length; i += BATCH_SIZE) {
    const batch = items.slice(i, i + BATCH_SIZE);
    console.log(`Processing batch ${Math.floor(i / BATCH_SIZE) + 1}/${Math.ceil(items.length / BATCH_SIZE)}: ${batch.length} items`);
    
    try {
      const batchResults = await batchProcessor(batch);
      results.push(...batchResults);
    } catch (error) {
      console.error(`Batch ${Math.floor(i / BATCH_SIZE) + 1} failed:`, error);
      // 添加空結果以保持索引一致
      results.push(...new Array(batch.length).fill(null));
    }
    
    // 如果不是最後一批，等待一段時間
    if (i + BATCH_SIZE < items.length) {
      console.log(`Waiting ${BATCH_DELAY}ms before next batch...`);
      await delay(BATCH_DELAY);
    }
  }
  
  return results;
};

// 圖表數據專用的批次處理（更快的併發）
const processChartBatches = async (items, batchProcessor) => {
  const results = [];
  const totalBatches = Math.ceil(items.length / CHART_BATCH_SIZE);
  
  for (let i = 0; i < items.length; i += CHART_BATCH_SIZE) {
    const batch = items.slice(i, i + CHART_BATCH_SIZE);
    const batchNumber = Math.floor(i / CHART_BATCH_SIZE) + 1;
    const requestCount = batch.length * 2; // 每個區塊需要 2 個請求
    console.log(`📊 圖表數據批次 ${batchNumber}/${totalBatches}: 併發處理 ${batch.length} 個區塊 (${requestCount} 個請求)`);
    
    try {
      const batchResults = await batchProcessor(batch);
      results.push(...batchResults);
      console.log(`✅ 批次 ${batchNumber} 完成，成功處理 ${batchResults.filter(r => r !== null).length}/${batch.length} 個區塊`);
    } catch (error) {
      console.error(`❌ 圖表數據批次 ${batchNumber} 失敗:`, error);
      // 添加空結果以保持索引一致
      results.push(...new Array(batch.length).fill(null));
    }
    
    // 如果不是最後一批，等待較短時間
    if (i + CHART_BATCH_SIZE < items.length) {
      console.log(`⏳ 頻率控制: 等待 ${CHART_BATCH_DELAY}ms 後處理下一批次 (保持 ≤14 req/sec)...`);
      await delay(CHART_BATCH_DELAY);
    }
  }
  
  return results;
};
  
export const ether_obol = {
  getObolOperatorClustersRegistry: async () => {
    try {
      console.log('Fetching active node operators count...');
      const activeNodeOperators = await ObolOperatorClustersRegistryContract.getActiveNodeOperatorsCount()
      const operatorCount = Number(activeNodeOperators) // 轉換 BigNumber 為普通數字
      console.log('Active node operators count:', operatorCount)
      
      if (operatorCount === 0) {
        console.log('No active node operators found')
        return []
      }
      
      // 創建操作者索引數組
      const operatorIndices = Array.from({ length: operatorCount }, (_, i) => i);
      
      // 批次處理器函數
      const batchProcessor = async (indices) => {
        const promises = indices.map(index => 
          ObolOperatorClustersRegistryContract.getNodeOperator(index, true).catch(error => {
            console.error(`Error fetching operator ${index}:`, error);
            return null;
          })
        );
        
        return Promise.all(promises);
      };
      
      console.log(`Starting to fetch ${operatorCount} operators in batches of ${BATCH_SIZE}...`);
      const nodeOperators = await processBatches(operatorIndices, batchProcessor);
      
      // 過濾掉失敗的調用
      const validOperators = nodeOperators.filter(operator => operator !== null);
      console.log(`Successfully fetched ${validOperators.length}/${operatorCount} operators`);
      
      return validOperators;
    } catch (error) {
      console.error('Error in getObolOperatorClustersRegistry:', error);
      throw error;
    }
  },
  getObolOperatorHistoryValidators: async (id, period) => {
    const startTime = Date.now();
    console.log(`🚀 開始載入操作者 ${id} 的 ${period} 歷史數據`);

    const PERIOD_CONFIG = {
      '1m': {days: 30, stepDays: 3},
      '1y': {days: 365, stepDays: 30},
      'all': {days: 730, stepDays: 50} // 'all' 設為 2 年，每 30 天一個數據點
    };

    const { days, stepDays } = PERIOD_CONFIG[period];
    const latestBlock = await provider.getBlockNumber();
    const latestBlkNum = BigInt(latestBlock);
    
    const startBlkNum = latestBlkNum - BigInt(Math.floor(days * BLOCKS_PER_DAY));
    const stepBlocks = BigInt(Math.floor(stepDays * BLOCKS_PER_DAY));
    
    // 生成所有需要查詢的區塊號
    const blocks = [];
    for (let b = startBlkNum; b <= latestBlkNum; b += stepBlocks) {
      blocks.push(Number(b));
    }
    
    const estimatedRequests = blocks.length * 2; // 每個區塊需要 2 個請求
    const estimatedTime = Math.ceil(blocks.length / CHART_BATCH_SIZE) * (CHART_BATCH_DELAY / 1000);
    console.log(`📋 需要查詢 ${blocks.length} 個區塊的數據，預估時間段範圍: ${Math.round(days)} 天`);
    console.log(`🚦 請求配置: 每批 ${CHART_BATCH_SIZE} 個區塊 (${CHART_BATCH_SIZE * 2} 個請求), 批次間隔 ${CHART_BATCH_DELAY}ms`);
    console.log(`⏱️ 預估總請求數: ${estimatedRequests} 個，預估完成時間: ${estimatedTime} 秒`);

    // 批次處理器函數 - 併發處理區塊數據
    const batchProcessor = async (blockBatch) => {
      const promises = blockBatch.map(async (blockNumber) => {
        const maxRetries = 2;
        let lastError;
        
        for (let retry = 0; retry <= maxRetries; retry++) {
          try {
            // 併發執行兩個請求：操作者數據和區塊信息
            const [validator, blockInfo] = await Promise.all([
              ObolOperatorClustersRegistryContract.getNodeOperator(id, true, {blockTag: blockNumber}),
              provider.getBlock(blockNumber)
            ]);

            return {
              blockNumber: blockNumber,
              timestamp: blockInfo.timestamp * 1000, // 轉換為毫秒
              date: new Date(blockInfo.timestamp * 1000).toISOString().split('T')[0],
              data: {
                active: validator[0],
                name: validator[1],
                rewardAddress: validator[2],
                totalVettedValidators: Number(validator[3]),
                totalExitedValidators: Number(validator[4]),
                totalAddedValidators: Number(validator[5]),
                totalDepositedValidators: Number(validator[6])
              }
            };
                      } catch (error) {
              lastError = error;
              if (retry < maxRetries) {
                // 檢查是否為頻率限制錯誤
                const isRateLimitError = error.message?.includes('request limit reached') || 
                                       error.code === 'UNKNOWN_ERROR' && error.error?.code === -32007;
                
                // 如果是頻率限制錯誤，延遲更長時間
                const retryDelay = isRateLimitError ? 2000 + (retry * 1000) : 500 * (retry + 1);
                
                await delay(retryDelay);
                console.warn(`🔄 區塊 ${blockNumber} 第 ${retry + 1} 次重試 (延遲 ${retryDelay}ms)...`);
              }
            }
        }
        
        console.warn(`❌ 區塊 ${blockNumber} 載入失敗 (${maxRetries + 1} 次嘗試):`, lastError?.message || 'Unknown error');
        return null; // 返回 null 表示失敗
      });

      return Promise.all(promises);
    };

    // 使用圖表專用的批次處理來併發載入數據（更快的併發處理）
    const results = await processChartBatches(blocks, batchProcessor);
    
    // 過濾掉失敗的請求並按時間戳排序
    const validators = results
      .filter(result => result !== null)
      .sort((a, b) => a.timestamp - b.timestamp);
    
    const endTime = Date.now();
    const duration = (endTime - startTime) / 1000;
    const successRate = ((validators.length / blocks.length) * 100).toFixed(1);
    const actualRequestCount = blocks.length * 2; // 總請求數
    const actualRequestRate = (actualRequestCount / duration).toFixed(1);
    
    console.log(`🎉 載入完成！成功 ${validators.length}/${blocks.length} 個數據點 (${successRate}%)`);
    console.log(`⚡ 總耗時: ${duration.toFixed(2)} 秒，平均每個區塊: ${(duration / blocks.length).toFixed(3)} 秒`);
    console.log(`🚦 實際請求頻率: ${actualRequestRate} req/sec (限制: 15 req/sec)`);
    
    return validators;
  },
  getObolOperatorSplitWallets: async (operatorAddress) => {
    const operatorContract = new ethers.Contract(operatorAddress, ObolOperatorABI, provider);
    const splitWallet = await operatorContract.splitWallet();
    return splitWallet;
  },
  getObolOperatorRewardshare: async (splitWallet) => {
    const url = `https://api.etherscan.io/api?module=account&action=txlistinternal&address=${splitWallet}&startblock=0&endblock=99999999&sort=desc&apikey=${EtherscanAPIKey}`;
    console.log('Fetching transactions for split wallet:', splitWallet);
    const transactions = await fetch(url).then(res => res.json()).then(data => data.result);
    console.log('Transactions:', transactions);
    const createTx = transactions.find(tx => tx.type == "create2")
    console.log('Create tx:', createTx);
    const createOriginTx = await etherscanProvider.getTransaction(createTx.hash);
    console.log('Create origin tx:', createOriginTx.data);
    const decodedData = SplitMainInterface.decodeFunctionData('createSplit', createOriginTx.data);
    console.log('Decoded data:', decodedData);
    if (decodedData) {
      const rewardAddress = decodedData[0];
      const rewardShare = decodedData[1];
      console.log('Reward address:', rewardAddress);
      console.log('Reward share:', rewardShare);
      return {rewardAddress,rewardShare};
    } else {
      throw new Error('Failed to decode createSplit function data');
    }
  },

  getObolPredictedRewardShare: async (operatorAddress) => {
    const operatorContract = new ethers.Contract(operatorAddress, ObolOperatorABI, provider);
    const splitWallet = await operatorContract.splitWallet();
    const rewardShare = await ether_obol.getObolOperatorRewardshare(splitWallet);
    return rewardShare;
  },

>>>>>>> 09f398e (update ValidatorHistory)
}