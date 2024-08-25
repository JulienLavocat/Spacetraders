import { useQuery } from "react-query";
import { API_URL } from "./api";

export interface Waypoint {
  id: string;
  type: string;
  faction: string;
  orbits: string | null;
  submittedBy: string | null;
  submittedOn: string | null;
  x: number;
  y: number;
  isUnderConstruction: boolean;
}

export function useWaypointsList(systemId: string) {
  return useQuery({
    queryKey: "waypoints-" + systemId,
    queryFn: async () => {
      const res = await fetch(`${API_URL}/systems/${systemId}`);
      const results = (await res.json()) as Waypoint[];

      const waypoints: Record<string, Waypoint> = {};
      const orbits: Record<string, string[]> = {};

      for (const w of results) {
        waypoints[w.id] = w;

        if (!w.orbits) orbits[w.id] = [];
        else orbits[w.orbits].push(w.id);
      }

      Object.values(orbits).forEach((e) => e.sort());

      return { waypoints, orbits };
    },
  });
}
