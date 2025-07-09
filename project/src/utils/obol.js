import { ethers } from "ethers";

// Operator Clusters Registry Contract
const ObolOperatorClustersRegistry = "0xaE7B191A31f627b4eB1d4DaC64eaB9976995b433";

// Obol Operator Clusters Registry Contract ABI
const ObolOperatorClustersRegistryProxyABI = [{"constant":true,"inputs":[],"name":"proxyType","outputs":[{"name":"proxyTypeId","type":"uint256"}],"payable":false,"stateMutability":"pure","type":"function"},{"constant":true,"inputs":[],"name":"isDepositable","outputs":[{"name":"","type":"bool"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[],"name":"implementation","outputs":[{"name":"","type":"address"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[],"name":"appId","outputs":[{"name":"","type":"bytes32"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[],"name":"kernel","outputs":[{"name":"","type":"address"}],"payable":false,"stateMutability":"view","type":"function"},{"inputs":[{"name":"_kernel","type":"address"},{"name":"_appId","type":"bytes32"},{"name":"_initializePayload","type":"bytes"}],"payable":false,"stateMutability":"nonpayable","type":"constructor"},{"payable":true,"stateMutability":"payable","type":"fallback"},{"anonymous":false,"inputs":[{"indexed":false,"name":"sender","type":"address"},{"indexed":false,"name":"value","type":"uint256"}],"name":"ProxyDeposit","type":"event"}]
const ObolOperatorClustersRegistryLogicABI = [{"constant":true,"inputs":[],"name":"getActiveNodeOperatorsCount","outputs":[{"name":"","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[{"name":"_nodeOperatorId","type":"uint256"},{"name":"_fullInfo","type":"bool"}],"name":"getNodeOperator","outputs":[{"name":"active","type":"bool"},{"name":"name","type":"string"},{"name":"rewardAddress","type":"address"},{"name":"totalVettedValidators","type":"uint64"},{"name":"totalExitedValidators","type":"uint64"},{"name":"totalAddedValidators","type":"uint64"},{"name":"totalDepositedValidators","type":"uint64"}],"payable":false,"stateMutability":"view","type":"function"}];
const fullABI = [...ObolOperatorClustersRegistryProxyABI, ...ObolOperatorClustersRegistryLogicABI];

const RPC_URL = "https://convincing-divine-slug.quiknode.pro/9e835fa839e702f4140f3a80b65672cfb3d5950b/";
const provider = new ethers.JsonRpcProvider(RPC_URL);
const contract = new ethers.Contract(ObolOperatorClustersRegistry, fullABI, provider);

// 請求節流配置
const BATCH_SIZE = 10; // 每批最多10個請求
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
      const activeNodeOperators = await contract.getActiveNodeOperatorsCount()
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
          contract.getNodeOperator(index, true).catch(error => {
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
  }
}