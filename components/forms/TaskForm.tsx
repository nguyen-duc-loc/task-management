"use client";

import { zodResolver } from "@hookform/resolvers/zod";
import { IconCalendarWeek } from "@tabler/icons-react";
import React from "react";
import { useForm } from "react-hook-form";
import { toast } from "sonner";

import { Button } from "@/components/ui/button";
import {
  FormField,
  FormItem,
  FormLabel,
  FormControl,
  FormMessage,
  Form,
  FormDescription,
} from "@/components/ui/form";
import { Input } from "@/components/ui/input";
import { TaskData, TaskSchema } from "@/validation";
import { Textarea } from "../ui/textarea";
import { Popover, PopoverContent, PopoverTrigger } from "../ui/popover";
import { cn } from "@/lib/utils";
import { format } from "date-fns";
import { ScrollArea, ScrollBar } from "../ui/scroll-area";
import { Calendar } from "../ui/calendar";
import { Plus, Check } from "lucide-react";
import Spinner from "../ui/spinner";

interface TaskFormProps {
  defaultValues: TaskData;
  onSubmit: (data: TaskData) => Promise<ActionResponse<unknown>>;
  formType: "CREATE" | "UPDATE";
}

const TaskForm = ({ defaultValues, onSubmit, formType }: TaskFormProps) => {
  const form = useForm({
    resolver: zodResolver(TaskSchema),
    defaultValues,
  });

  const handleSubmit = async (data: TaskData) => {
    const response = await onSubmit(data);
    if (!response.success) {
      toast.error(
        `${response.error?.[0].toUpperCase()}${response.error?.slice(1)}`
      );
    } else {
      toast.success(
        `${formType[0]}${formType.slice(1).toLowerCase()}d task successfully`
      );
      form.reset();
    }
  };

  const handleDateSelect = (date?: Date) => {
    if (date) {
      form.setValue("deadline", date);
    }
  };

  const handleTimeChange = (type: "hour" | "minute", value: string) => {
    const currentDate = form.getValues("deadline") || new Date();
    let newDate = new Date(currentDate);

    if (type === "hour") {
      const hour = parseInt(value, 10);
      newDate.setHours(hour);
    } else if (type === "minute") {
      newDate.setMinutes(parseInt(value, 10));
    }

    form.setValue("deadline", newDate);
  };

  const isSubmitting = form.formState.isSubmitting;

  return (
    <Form {...form}>
      <form
        onSubmit={form.handleSubmit(handleSubmit)}
        className="sm:mx-12 space-y-8"
      >
        <FormField
          control={form.control}
          name="title"
          render={({ field }) => (
            <FormItem>
              <FormLabel>
                Title <span className="text-primary">*</span>
              </FormLabel>
              <FormControl>
                <Input
                  {...field}
                  disabled={isSubmitting}
                  className="h-10"
                  placeholder="Your task title"
                />
              </FormControl>
              <FormDescription>
                A short, clear name for the task.
              </FormDescription>
              <FormMessage />
            </FormItem>
          )}
        />

        <FormField
          control={form.control}
          name="description"
          render={({ field }) => (
            <FormItem>
              <FormLabel>Description</FormLabel>
              <FormControl>
                <Textarea
                  {...field}
                  disabled={isSubmitting}
                  className="h-40"
                  placeholder="Add details about the task, steps to complete, or any notes"
                />
              </FormControl>
              <FormDescription>
                Provide more context so you know exactly what's involved.
              </FormDescription>
              <FormMessage />
            </FormItem>
          )}
        />

        <FormField
          control={form.control}
          name="deadline"
          render={({ field }) => (
            <FormItem>
              <FormLabel>
                Deadline <span className="text-primary">*</span>
              </FormLabel>
              <Popover>
                <PopoverTrigger asChild>
                  <FormControl>
                    <Button
                      disabled={isSubmitting}
                      variant={"outline"}
                      className={cn(
                        "w-fit pl-3 text-left font-normal",
                        !field.value && "text-muted-foreground"
                      )}
                    >
                      {field.value ? (
                        format(field.value, "yyyy/MM/dd HH:mm")
                      ) : (
                        <span>YYYY/MM/DD HH:mm</span>
                      )}
                      <IconCalendarWeek className="ml-auto h-4 w-4 opacity-50" />
                    </Button>
                  </FormControl>
                </PopoverTrigger>
                <PopoverContent className="w-auto p-0">
                  <div className="sm:flex">
                    <Calendar
                      mode="single"
                      selected={field.value}
                      onSelect={handleDateSelect}
                      initialFocus
                      disabled={{ before: new Date() }}
                    />
                    <div className="flex flex-col sm:flex-row sm:h-[300px] divide-y sm:divide-y-0 sm:divide-x">
                      <ScrollArea className="w-64 sm:w-auto">
                        <div className="flex sm:flex-col p-2">
                          {Array.from({ length: 24 }, (_, i) => i)
                            .reverse()
                            .map((hour) => (
                              <Button
                                key={hour}
                                size="icon"
                                variant={
                                  field.value && field.value.getHours() === hour
                                    ? "default"
                                    : "ghost"
                                }
                                className="sm:w-full shrink-0 aspect-square"
                                onClick={() =>
                                  handleTimeChange("hour", hour.toString())
                                }
                              >
                                {hour}
                              </Button>
                            ))}
                        </div>
                        <ScrollBar
                          orientation="horizontal"
                          className="sm:hidden"
                        />
                      </ScrollArea>
                      <ScrollArea className="w-64 sm:w-auto">
                        <div className="flex sm:flex-col p-2">
                          {Array.from({ length: 12 }, (_, i) => i * 5).map(
                            (minute) => (
                              <Button
                                key={minute}
                                size="icon"
                                variant={
                                  field.value &&
                                  field.value.getMinutes() === minute
                                    ? "default"
                                    : "ghost"
                                }
                                className="sm:w-full shrink-0 aspect-square"
                                onClick={() =>
                                  handleTimeChange("minute", minute.toString())
                                }
                              >
                                {minute.toString().padStart(2, "0")}
                              </Button>
                            )
                          )}
                        </div>
                        <ScrollBar
                          orientation="horizontal"
                          className="sm:hidden"
                        />
                      </ScrollArea>
                    </div>
                  </div>
                </PopoverContent>
              </Popover>
              <FormDescription>
                Choose when this task should be completed by.
              </FormDescription>
              <FormMessage />
            </FormItem>
          )}
        />

        <div className="flex justify-end">
          <Button disabled={isSubmitting} className="ml-auto w-fit">
            {isSubmitting ? (
              <Spinner />
            ) : formType === "CREATE" ? (
              <Plus />
            ) : (
              <Check />
            )}
            {formType[0].toUpperCase() + formType.slice(1).toLowerCase()}
          </Button>
        </div>
      </form>
    </Form>
  );
};

export default TaskForm;
