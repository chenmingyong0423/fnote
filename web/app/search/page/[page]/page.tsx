import { searchMetadata, renderSearch } from "@/src/components/SearchResultsPage";
import type { ListSearchParams } from "@/src/utils/pagination";

type Props = {
  params: Promise<{ page: string }>;
  searchParams: Promise<ListSearchParams>;
};

export async function generateMetadata({ params, searchParams }: Props) {
  return searchMetadata((await params).page, await searchParams);
}

export default async function Page({ params, searchParams }: Props) {
  return renderSearch((await params).page, await searchParams);
}
