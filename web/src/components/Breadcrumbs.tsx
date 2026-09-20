import Link from "next/link";
import { buildBreadcrumbJsonLd, serializeJsonLd, type BreadcrumbItem } from "@/src/utils/seo";

export default function Breadcrumbs({ items }: { items: BreadcrumbItem[] }) {
  return (
    <>
      <script type="application/ld+json" dangerouslySetInnerHTML={{ __html: serializeJsonLd(buildBreadcrumbJsonLd(items)) }} />
      <nav aria-label="面包屑导航" className="mb-4 text-sm text-gray-500 dark:text-gray-400">
        <ol className="flex flex-wrap items-center gap-x-2 gap-y-1">
          {items.map((item, index) => (
            <li key={item.pathname} className="min-w-0 break-words">
              {index > 0 && <span aria-hidden="true" className="mr-2">/</span>}
              {index === items.length - 1
                ? <span aria-current="page">{item.name}</span>
                : <Link href={item.pathname} className="hover:underline">{item.name}</Link>}
            </li>
          ))}
        </ol>
      </nav>
    </>
  );
}
