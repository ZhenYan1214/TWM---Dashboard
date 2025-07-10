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
}