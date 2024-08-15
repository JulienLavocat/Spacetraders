import { useMatches } from "react-router-dom";

export interface BreadcrumbData {
  link: string;
  name: string;
}

export function useBreadcrumbs(): BreadcrumbData[] {
  const matches = useMatches();
  return matches
    .filter((match) => (match.handle as any)?.crumb)
    .map((match) => (match.handle as any)?.crumb(match.params));
}
