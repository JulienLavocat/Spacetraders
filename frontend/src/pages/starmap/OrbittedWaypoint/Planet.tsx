import { Waypoint } from "@/api/waypoints";

export function Planet({
  waypoint,
  orbits,
}: {
  waypoint: Waypoint;
  orbits: Waypoint[];
}) {
  return (
    <g>
      <circle
        cx={waypoint.x}
        cy={waypoint.y}
        r={100}
        fill="teal"
        className="cursor-pointer"
      ></circle>
    </g>
  );
}
