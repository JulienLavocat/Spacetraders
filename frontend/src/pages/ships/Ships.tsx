import { Ship, useListShips } from "@/api/api";
import { DataTable, DataTableSortedHeader } from "@/components/ui/data-table";
import { ColumnDef } from "@tanstack/react-table";
import {
  PayloadColumn,
  RouteColumn,
  TradeRouteColumn,
  WaypointColumn,
} from "./Columns";

const columns: ColumnDef<Ship>[] = [
  {
    header: ({ column }) => <DataTableSortedHeader name="Id" column={column} />,
    accessorKey: "id",
  },
  {
    header: ({ column }) => (
      <DataTableSortedHeader name="System" column={column} />
    ),
    accessorKey: "system",
  },
  {
    header: "Waypoint",
    cell: ({ row }) => <WaypointColumn ship={row.original} />,
  },
  {
    header: "Payload",
    cell: ({ row }) => <PayloadColumn ship={row.original} />,
  },
  {
    header: "Route",
    cell: ({ row }) => <RouteColumn ship={row.original} />,
  },
  {
    header: "Trade route",
    cell: ({ row }) => <TradeRouteColumn ship={row.original} />,
  },
];

export function Ships() {
  const { data, isError, isLoading } = useListShips();

  if (isLoading) return <p>Loading...</p>;
  if (isError) return <p>An error occured</p>;

  return (
    <div className="font-mono">
      <DataTable
        columns={columns}
        data={data ?? []}
        filterColumn="id"
        filterPlaceholder="Filter ships by ID"
      />
    </div>
  );
}
