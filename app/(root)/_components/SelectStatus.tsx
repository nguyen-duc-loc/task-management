"use client";

import {
  Select,
  SelectTrigger,
  SelectValue,
  SelectContent,
  SelectItem,
} from "@/components/ui/select";
import { formUrlQuery, removeKeysFromQuery } from "@/lib/url";
import { useRouter, useSearchParams } from "next/navigation";
import router from "next/router";
import React from "react";

type Status = "all" | "todo" | "completed";
const options: {
  title: string;
  value: Status;
}[] = [
  {
    title: "All",
    value: "all",
  },
  {
    title: "To do",
    value: "todo",
  },
  {
    title: "Completed",
    value: "completed",
  },
];

const SelectStatus = () => {
  const searchParams = useSearchParams();
  const router = useRouter();
  const completed = searchParams.get("completed");

  const handleValueChange = (value: Status) => {
    let newUrl: string;

    if (value === "all") {
      newUrl = removeKeysFromQuery({
        params: searchParams.toString(),
        keysToRemove: ["completed"],
      });
    } else {
      newUrl = formUrlQuery({
        params: searchParams.toString(),
        key: "completed",
        value: value === "completed" ? "true" : "false",
      });
    }

    newUrl = removeKeysFromQuery({
      params: newUrl.slice(2),
      keysToRemove: ["page"],
    });

    router.push(newUrl, { scroll: false });
  };

  return (
    <Select
      defaultValue={
        completed === "true"
          ? "completed"
          : completed === "false"
          ? "todo"
          : "all"
      }
      onValueChange={handleValueChange}
    >
      <SelectTrigger className="w-32">
        <SelectValue placeholder="Status" />
      </SelectTrigger>
      <SelectContent>
        {options.map(({ title, value }) => (
          <SelectItem key={value} value={value}>
            {title}
          </SelectItem>
        ))}
      </SelectContent>
    </Select>
  );
};

export default SelectStatus;
