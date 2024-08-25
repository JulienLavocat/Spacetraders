import { useQuery } from "react-query";
import { API_URL } from "./api";

export interface System {
  id: string;
  type: string;
  x: number;
  y: number;
}

export function useListSystems() {
  return useQuery<System[]>({
    queryKey: "listSystems",
    queryFn: async () => {
      const res = await fetch(`${API_URL}/systems`);
      return res.json();
    },
  });
}
