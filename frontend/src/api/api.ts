import { useQuery } from "react-query";

export interface Route {
  to: string;
  fuel: number;
  aggFuel: number;
}

export interface TradeRoute {
  product: string;
  buyAt: string;
  sellAt: string;
  maxAmount: number;
  buyPrice: number;
  sellPrice: number;
  estimatedProfits: number;
  estimatedProfitsPerUnit: number;
  fuelCost: number;
}

export interface Ship {
  id: string;
  arrivalAt: string;
  departedAt: string;
  cooldown: string;
  cargo: Record<string, number>;
  tradeRoute: TradeRoute;
  route: Route[];
  status: string;
  origin: string;
  destination: string;
  system: string;
  waypoint: string;
  maxFuel: number;
  currentFuel: number;
  maxCargo: number;
  currentCargo: number;
  isCargoFull: boolean;
}

export const API_URL = "http://localhost:8080";

export function createQueryParams(params: Record<string, any>): string {
  const filtered = Object.fromEntries(
    Object.entries(params)
      .map(([key, value]) => (!value ? null : [key, value]))
      .filter((e) => !!e),
  );

  return new URLSearchParams(filtered).toString();
}

export function useListShips() {
  return useQuery<Ship[]>({
    queryKey: "listShips",
    refetchInterval: 2000,
    queryFn: async () => {
      const res = await fetch(`${API_URL}/ships`);
      return res.json();
    },
  });
}
