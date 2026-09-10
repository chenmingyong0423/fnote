import { searchMetadata, renderSearch } from "@/src/components/SearchResultsPage";
import type { ListSearchParams } from "@/src/utils/pagination";

type Props = {
  searchParams: Promise<ListSearchParams>;
};

export async function generateMetadata({ searchParams }: Props) {
  return searchMetadata(undefined, await searchParams);
}

export default async function Page({ searchParams }: Props) {
  return renderSearch(undefined, await searchParams);
}
