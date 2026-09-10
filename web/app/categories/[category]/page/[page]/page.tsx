import { archiveMetadata, renderArchive } from "@/src/components/ArchivePage";
import type { ListSearchParams } from "@/src/utils/pagination";

type Props = {
  params: Promise<{ category: string; page: string; }>;
  searchParams: Promise<ListSearchParams>;
};

export async function generateMetadata({ params, searchParams }: Props) {
  const route = await params;
  return archiveMetadata("categories", route.category, route.page, await searchParams);
}

export default async function Page({ params, searchParams }: Props) {
  const route = await params;
  return renderArchive("categories", route.category, route.page, await searchParams);
}
