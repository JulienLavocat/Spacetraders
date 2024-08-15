import { useQuery } from "react-query";
import { API_URL, createQueryParams } from "./api";

export interface Transaction {
  id: number;
  waypoint: string;
  product: string;
  type: string;
  ship: string;
  correlationId: string;
  totalPrice: number;
  agentBalance: number;
  amount: number;
  pricePerUnit: number;
  timestamp: string;
}

export interface ListTransactionsResult {
  transactions: Transaction[];
  revenue: number;
  expenses: number;
  fuelExpenses: number;
}

export interface ListTransactionsParams {
  page: number;
  limit: number;
  correlationId?: string;
}

export function useListTransactions(params: ListTransactionsParams) {
  const queryParams = createQueryParams(params);
  return useQuery<ListTransactionsResult>({
    queryKey: "listTransactions-" + params.correlationId,
    refetchInterval: 10000,
    cacheTime: 30000,
    queryFn: async () => {
      const res = await fetch(
        `${API_URL}/transactions?${queryParams.toString()}`,
      );
      return res.json();
    },
  });
}
