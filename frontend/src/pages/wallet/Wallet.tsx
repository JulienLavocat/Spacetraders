import { useParams } from "react-router-dom";
import { TransactionsTable } from "./transactions/Transactions";
import { countReset } from "console";

export function Wallet() {
  const correlationId = useParams().id;
  const limit = correlationId ? 1000 : 20;
  return (
    <TransactionsTable
      title={
        correlationId
          ? "Transactions related to " + correlationId
          : `Last ${limit} transactions`
      }
      showAgentBalance={!correlationId}
      params={{
        page: 1,
        limit,
        correlationId,
      }}
    />
  );
}
