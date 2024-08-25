import {
  ListTransactionsParams,
  Transaction,
  useListTransactions,
} from "@/api/transactions";
import {
  Card,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { DataTable } from "@/components/ui/data-table";
import { TypographyH2 } from "@/components/ui/typography";
import { cn } from "@/lib/utils";
import { formatCurrency } from "@/utils/format-currency";
import { createDurationFormatter } from "@/utils/format-duration";
import { ColumnDef } from "@tanstack/react-table";
import { CorrelationIdColumn } from "../Columns";

const formatDuration = createDurationFormatter("narrow");

const columns: ColumnDef<Transaction>[] = [
  {
    header: "Id",
    accessorKey: "id",
  },
  {
    header: "Ship",
    accessorKey: "ship",
  },
  {
    header: "Waypoint",
    accessorKey: "waypoint",
  },
  {
    header: "Product",
    accessorKey: "product",
  },
  {
    header: "Amount",
    accessorKey: "amount",
  },
  {
    header: "Type",
    accessorKey: "type",
  },
  {
    header: "Price",
    accessorFn: (row) =>
      `${formatCurrency(row.totalPrice)} (${formatCurrency(row.pricePerUnit)}/u)`,
  },
  {
    header: "Date",
    accessorFn: (row) => {
      const date = new Date(row.timestamp);
      const ellapsed = new Date().getTime() - date.getTime();

      if (ellapsed < 24 * 3600 * 1000) {
        return formatDuration(ellapsed) + " ago";
      }

      return Intl.DateTimeFormat("fr-FR", {
        timeStyle: "medium",
        dateStyle: "short",
        timeZone: "Europe/Paris",
      }).format(new Date(row.timestamp));
    },
  },
  {
    header: "Correlation ID",
    cell: ({ row }) => <CorrelationIdColumn tx={row.original} />,
  },
];

function TransactionStatsCard({
  title,
  amount,
  color,
}: {
  title: string;
  amount: number;
  color?: string;
}) {
  const isPositive = amount > 0;
  return (
    <Card>
      <CardHeader className="pb-2">
        <CardDescription>{title}</CardDescription>
        <CardTitle
          className={cn(
            "text-4xl",
            isPositive && !color ? "text-green-600" : "text-red-600",
            color,
          )}
        >
          {formatCurrency(amount)}
        </CardTitle>
      </CardHeader>
    </Card>
  );
}

export function TransactionsTable({
  params,
  showAgentBalance,
  title,
}: {
  params: ListTransactionsParams;
  showAgentBalance?: boolean;
  title?: string;
}) {
  const { data, isError, isLoading } = useListTransactions(params);

  if (isLoading) return <p>Loading...</p>;
  if (isError) return <p>An error occured</p>;

  if (!data) return;

  const agentBalance = data.transactions[0].agentBalance;

  return (
    <>
      <TypographyH2>{title}</TypographyH2>
      <div
        className={cn(
          "grid gap-16",
          showAgentBalance ? "grid-cols-5" : "grid-cols-4",
        )}
      >
        {showAgentBalance && (
          <TransactionStatsCard title="Agent balance" amount={agentBalance} />
        )}
        <TransactionStatsCard
          title="Revenue"
          color="text-green-600"
          amount={data.revenue}
        />
        <TransactionStatsCard
          title="Expenses"
          color="text-red-600"
          amount={data.expenses}
        />
        <TransactionStatsCard
          title="Fuel expenses"
          color="text-yellow-600"
          amount={data.fuelExpenses}
        />
        <TransactionStatsCard
          title="Profit"
          amount={data.revenue - data.expenses}
        />
      </div>
      <div className="font-mono">
        <DataTable columns={columns} data={data?.transactions ?? []} />
      </div>
    </>
  );
}
