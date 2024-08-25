import { useWaypointsList } from "@/api/waypoints";
import { useParams } from "react-router-dom";
import { TransformComponent, TransformWrapper } from "react-zoom-pan-pinch";
import { OrbittedWaypoint } from "./OrbittedWaypoint/OrbittedWaypoint";

export function SystemMap() {
  const systemId = useParams().id as string;
  const { data, isLoading, isError } = useWaypointsList(systemId);

  if (isLoading) return <p>Loading...</p>;
  if (isError || !data)
    return <p>An error occured while loading the systems</p>;

  const { waypoints, orbits } = data;

  return (
    <div className="bg-slate-900 max-w-full overflow-hidden font-mono">
      <TransformWrapper
        limitToBounds={false}
        // initialScale={15}
        // minScale={1}
        maxScale={15}
        initialScale={2}
        initialPositionX={0}
        initialPositionY={0}
      >
        <TransformComponent>
          <svg
            width="100vw"
            height="calc(100vh - 4.5rem)"
            viewBox="-789.0177437802017 -789.0177437802017 1578.0354875604035 1578.0354875604035"
          >
            {Object.entries(orbits).map(([waypoint, orbits]) => (
              <OrbittedWaypoint
                waypoint={waypoints[waypoint]}
                orbits={orbits.map((wp) => waypoints[wp])}
              />
            ))}
          </svg>
        </TransformComponent>
      </TransformWrapper>
    </div>
  );
}
