import { Transaction } from "@/api/transactions";
import { Button } from "@/components/ui/button";
import { Link } from "react-router-dom";

export function CorrelationIdColumn({ tx }: { tx: Transaction }) {
  return (
    <Button variant="link" className="pl-0">
      <Link to={`/wallet/transactions/${tx.correlationId}`}>
        {tx.correlationId}
      </Link>
    </Button>
  );
}
