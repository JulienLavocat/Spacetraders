import { Tooltip, TooltipTrigger } from "@/components/ui/tooltip";
import { TooltipContent } from "@radix-ui/react-tooltip";
import dynamicIconImports from "lucide-react/dynamicIconImports";
import { lazy, Suspense } from "react";
import { Link } from "react-router-dom";

export interface NavLinkProps {
  icon: keyof typeof dynamicIconImports;
  text: string;
  link: string;
  active?: boolean;
}

const fallback = <div style={{ background: "#ddd" }} className="h-5 w-5" />;

export function NavLink({ icon, text, link, active }: NavLinkProps) {
  const LucideIcon = lazy(dynamicIconImports[icon]);
  return (
    <Tooltip>
      <TooltipTrigger asChild>
        <Link
          to={link}
          className={`flex h-9 w-9 items-center justify-center rounded-lg text-muted-foreground transition-colors hover:text-foreground md:h-8 md:w-8 ${active ? "bg-primary" : ""}`}
        >
          <Suspense fallback={fallback}>
            <LucideIcon className="h-5 w-5" />
          </Suspense>
          <span className="sr-only">{text}</span>
        </Link>
      </TooltipTrigger>
      <TooltipContent side="right">{text}</TooltipContent>
    </Tooltip>
  );
}
