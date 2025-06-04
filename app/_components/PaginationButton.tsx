"use client";

import { ChevronLeft, ChevronRight, MoreHorizontal } from "lucide-react";
import { useRouter, useSearchParams } from "next/navigation";
import React from "react";

import { Button } from "@/components/ui/button";
import { formUrlQuery, removeKeysFromQuery } from "@/lib/url";

const PaginationButton = ({
  page,
  isActive,
  disabled,
}: {
  page: number | "previous" | "next" | "ellipsis";
  isActive?: boolean;
  disabled?: boolean;
}) => {
  const searchParams = useSearchParams();
  const router = useRouter();
  const currentPage = Number(searchParams.get("page")) || 1;

  const handleClick = () => {
    if (isActive || page === "ellipsis") return;

    let newUrl: string;

    if (page === 1) {
      newUrl = removeKeysFromQuery({
        params: searchParams.toString(),
        keysToRemove: ["page"],
      });
    } else if (page === "previous") {
      if (currentPage === 2) {
        newUrl = removeKeysFromQuery({
          params: searchParams.toString(),
          keysToRemove: ["page"],
        });
      } else {
        newUrl = formUrlQuery({
          params: searchParams.toString(),
          key: "page",
          value: String(currentPage - 1),
        });
      }
    } else if (page === "next") {
      newUrl = formUrlQuery({
        params: searchParams.toString(),
        key: "page",
        value: String(currentPage + 1),
      });
    } else {
      newUrl = formUrlQuery({
        params: searchParams.toString(),
        key: "page",
        value: String(page),
      });
    }

    router.push(newUrl, { scroll: false });
  };

  return (
    <Button
      disabled={disabled}
      variant={
        isActive || page === "previous" || page === "next" ? "outline" : "ghost"
      }
      onClick={handleClick}
      className={`capitalize ${isActive && "font-semibold"} ${
        (page === "previous" || page === "next") && "p-3"
      } ${
        (page === "ellipsis" || typeof page === "number") &&
        "max-[500px]:hidden"
      }`}
    >
      {page === "previous" && <ChevronLeft />}
      {page === "ellipsis" ? (
        <MoreHorizontal />
      ) : typeof page === "number" ? (
        page
      ) : (
        ""
      )}
      {page === "next" && <ChevronRight />}
    </Button>
  );
};

export default PaginationButton;
