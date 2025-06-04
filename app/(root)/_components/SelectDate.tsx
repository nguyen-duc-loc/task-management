"use client";

import { Button } from "@/components/ui/button";
import { Calendar } from "@/components/ui/calendar";
import {
  Popover,
  PopoverTrigger,
  PopoverContent,
} from "@/components/ui/popover";
import { formUrlQuery, removeKeysFromQuery } from "@/lib/url";
import { cn } from "@/lib/utils";
import { format, isValid } from "date-fns";
import { CalendarIcon } from "lucide-react";
import { useRouter, useSearchParams } from "next/navigation";
import React from "react";
import { DateRange } from "react-day-picker";

const SelectDate = () => {
  const searchParams = useSearchParams();
  const router = useRouter();
  const from = searchParams.get("from") || "";
  const fromDate = isValid(new Date(from)) ? new Date(from) : undefined;
  const to = searchParams.get("to") || "";
  const toDate = isValid(new Date(to)) ? new Date(to) : undefined;
  const date: DateRange = {
    from: fromDate,
    to: toDate,
  };

  return (
    <div className="grid gap-2">
      <Popover>
        <PopoverTrigger asChild>
          <Button
            id="date"
            variant={"outline"}
            className={cn(
              "w-[260px] justify-start text-left font-normal",
              !date && "text-muted-foreground"
            )}
          >
            <CalendarIcon className="mr-2 h-4 w-4" />
            {date?.from ? (
              date.to ? (
                <>
                  {format(date.from, "LLL dd, y")} -{" "}
                  {format(date.to, "LLL dd, y")}
                </>
              ) : (
                `${format(date.from, "LLL dd, y")} -`
              )
            ) : (
              <span>Pick a date</span>
            )}
          </Button>
        </PopoverTrigger>
        <PopoverContent className="w-auto p-0" align="end">
          <Calendar
            mode="range"
            numberOfMonths={2}
            defaultMonth={date?.from}
            selected={date}
            onSelect={(date) => {
              let newUrl: string;

              if (!date) {
                newUrl = removeKeysFromQuery({
                  params: searchParams.toString(),
                  keysToRemove: ["from"],
                });
                newUrl = removeKeysFromQuery({
                  params: newUrl.slice(2),
                  keysToRemove: ["to"],
                });
              } else {
                if (date.from) {
                  newUrl = formUrlQuery({
                    params: searchParams.toString(),
                    key: "from",
                    value: format(date.from, "yyyy-MM-dd"),
                  });
                  if (date.to) {
                    newUrl = formUrlQuery({
                      params: newUrl.slice(2),
                      key: "to",
                      value: format(date.to, "yyyy-MM-dd"),
                    });
                  } else {
                    newUrl = removeKeysFromQuery({
                      params: newUrl.slice(2),
                      keysToRemove: ["to"],
                    });
                  }
                } else {
                  newUrl = removeKeysFromQuery({
                    params: searchParams.toString(),
                    keysToRemove: ["from"],
                  });
                  if (date.to) {
                    newUrl = formUrlQuery({
                      params: newUrl.slice(2),
                      key: "to",
                      value: format(date.to, "yyyy-MM-dd"),
                    });
                  } else {
                    newUrl = removeKeysFromQuery({
                      params: newUrl.slice(2),
                      keysToRemove: ["to"],
                    });
                  }
                }
              }

              newUrl = removeKeysFromQuery({
                params: newUrl.slice(2),
                keysToRemove: ["page"],
              });

              router.push(newUrl, { scroll: false });
            }}
          />
        </PopoverContent>
      </Popover>
    </div>
  );
};

export default SelectDate;
