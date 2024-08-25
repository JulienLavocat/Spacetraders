import { useListSystems } from "@/api/systems";
import L, { FeatureGroup, Map } from "leaflet";
import "leaflet/dist/leaflet.css";
import { createRef, useEffect, useState } from "react";
import {
  CircleMarker,
  FeatureGroup as FeatureGroupComponent,
  MapContainer,
  Tooltip,
} from "react-leaflet";
import { useNavigate } from "react-router-dom";

function SystemMarker({ x, y, name }: { x: number; y: number; name: string }) {
  const navigate = useNavigate();

  return (
    <CircleMarker
      center={[x, y]}
      radius={2}
      interactive={true}
      eventHandlers={{ click: () => navigate("/starmap/" + name) }}
    >
      <Tooltip direction={"top"} permanent={name === "X1-NT44"}>
        {name}
      </Tooltip>
    </CircleMarker>
  );
}

export function SystemsMap() {
  const { data, isLoading, isError } = useListSystems();
  const systemsMarkersRef = createRef<FeatureGroup<any>>();
  const [map, setMap] = useState<Map | null>(null);

  useEffect(() => {
    if (map != null && !isLoading && systemsMarkersRef.current) {
      map.fitBounds(systemsMarkersRef.current.getBounds());
    }
  }, [isLoading, map, systemsMarkersRef]);

  if (isLoading) return <p>Loading...</p>;
  if (isError) return <p>An error occured while loading the systems</p>;

  return (
    <MapContainer
      crs={L.CRS.Simple}
      preferCanvas={true}
      minZoom={-8}
      className="h-full !bg-background"
      ref={setMap}
      attributionControl={false}
    >
      <FeatureGroupComponent interactive={true} ref={systemsMarkersRef}>
        {data?.map((system) => (
          <SystemMarker
            key={system.id}
            x={system.x}
            y={system.y}
            name={system.id}
          />
        ))}
      </FeatureGroupComponent>
    </MapContainer>
  );
}
