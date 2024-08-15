import { Ship } from "@/api/api";
import { Progress } from "@/components/ui/progress";
import { formatCurrency } from "@/utils/format-currency";
import { createDurationFormatter } from "@/utils/format-duration";
import { AtSign, MoveRight } from "lucide-react";

function getWaypointId(fullWaypoint: string): string {
  return fullWaypoint.split("-")[2];
}

export function WaypointColumn({ ship }: { ship: Ship }) {
  const remaining = Math.round(
    new Date(ship.arrivalAt).getTime() - new Date().getTime(),
  );
  return (
    <>
      <div className="flex gap-2 items-center">
        <span>{getWaypointId(ship.origin)}</span>
        <MoveRight className="w-3" />
        <span>{getWaypointId(ship.destination)}</span>
        <span>{createDurationFormatter("narrow")(remaining)}</span>
      </div>
      <span className="text-neutral-400">{ship.status}</span>
    </>
  );
}
export function PayloadColumn({ ship }: { ship: Ship }) {
  return (
    <div className="grid grid-cols-payload grid-rows-2">
      <div className="flex gap-2 items-center">
        <span>F:</span>
        <Progress
          indicatorClassName="bg-yellow-600"
          className="h-2 mr-2"
          value={(ship.currentFuel / ship.maxFuel) * 100}
        />
      </div>
      {`(${ship.currentFuel}/${ship.maxFuel})`}
      <div className="flex gap-2 items-center">
        <span>C:</span>
        <Progress
          className="h-2 mr-2"
          value={(ship.currentCargo / ship.maxCargo) * 100}
        />
      </div>
      {`(${ship.currentCargo}/${ship.maxCargo})`}
    </div>
  );
}
export function RouteColumn({ ship }: { ship: Ship }) {
  if (!ship.route) {
    return null;
  }
  return (
    <div className="flex gap-2 items-center">
      {ship.route
        .map((route, i) => [
          <span>{getWaypointId(route.to)}</span>,
          i < ship.route.length - 1 ? <MoveRight className="w-3" /> : null,
        ])
        .flat()}
    </div>
  );
}
export function TradeRouteColumn({ ship }: { ship: Ship }) {
  const route = ship.tradeRoute;

  if (!route) {
    return null;
  }

  return (
    <>
      <div className="flex gap-2 items-center">
        <span>{getWaypointId(route.buyAt)}</span>
        <MoveRight className="w-3" />
        <span>{getWaypointId(route.sellAt)}</span>
      </div>
      <div className="flex gap-2 items-center">
        <p className="text-neutral-400">{route.maxAmount}</p>
        <p className="text-neutral-400">{route.product}</p>
        <AtSign className="w-3" />
        <span className="text-green-500">
          {formatCurrency(route.estimatedProfits)}
        </span>
        <span className="text-green-500">
          ({formatCurrency(route.estimatedProfitsPerUnit)}/unit)
        </span>
      </div>
    </>
  );
}
