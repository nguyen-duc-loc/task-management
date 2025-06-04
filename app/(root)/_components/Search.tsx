"use client";

import { Search as SearchIcon } from "lucide-react";
import { useSearchParams, useRouter, usePathname } from "next/navigation";
import React, { useEffect, useState } from "react";

import { Input } from "@/components/ui/input";
import { formUrlQuery, removeKeysFromQuery } from "@/lib/url";

const Search = () => {
  const pathname = usePathname();
  const searchParams = useSearchParams();
  const router = useRouter();
  const query = searchParams.get("search") || "";

  const [searchQuery, setSearchQuery] = useState(query);

  useEffect(() => {
    const delayDebounceFn = setTimeout(() => {
      let newUrl: string;

      if (searchQuery) {
        newUrl = formUrlQuery({
          params: searchParams.toString(),
          key: "search",
          value: searchQuery,
        });
      } else {
        newUrl = removeKeysFromQuery({
          params: searchParams.toString(),
          keysToRemove: ["search"],
        });
      }

      newUrl = removeKeysFromQuery({
        params: newUrl.slice(2),
        keysToRemove: ["page"],
      });

      router.push(newUrl, { scroll: false });
    }, 500);

    return () => clearTimeout(delayDebounceFn);
  }, [searchQuery, router, pathname]);

  return (
    <div className="flex min-h-10 max-w-96 grow items-center gap-2 rounded-[10px] border bg-background px-4">
      <SearchIcon className="cursor-pointer stroke-primary size-5" />
      <Input
        type="text"
        placeholder="Search task..."
        value={searchQuery}
        onChange={(e) => {
          setSearchQuery(e.target.value);
        }}
        className="border-none shadow-none outline-none focus-visible:ring-0 focus-visible:ring-transparent focus-visible:ring-offset-0"
      />
    </div>
  );
};

export default Search;
