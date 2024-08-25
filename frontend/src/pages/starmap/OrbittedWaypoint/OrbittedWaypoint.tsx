import { Waypoint } from "@/api/waypoints";
import { Planet } from "./Planet";

export function OrbittedWaypoint({
  waypoint,
  orbits,
}: {
  waypoint: Waypoint;
  orbits: Waypoint[];
}) {
  switch (waypoint.type) {
    case "PLANET":
      return <Planet waypoint={waypoint} orbits={orbits} />;
    default:
      return null;
  }
}
