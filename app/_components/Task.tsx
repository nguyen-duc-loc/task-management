import StatusBadge from "@/components/StatusBadge";
import {
  Card,
  CardAction,
  CardContent,
  CardFooter,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuGroup,
  DropdownMenuItem,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { Separator } from "@/components/ui/separator";
import ROUTES from "@/constants/routes";
import {
  IconCalendarWeek,
  IconDots,
  IconPencil,
  IconTrash,
} from "@tabler/icons-react";
import { format } from "date-fns";
import Link from "next/link";
import React from "react";
import MarkAsDoneBtn from "./MarkAsDoneBtn";
import DeleteBtn from "./DeleteBtn";

interface TaskProps {
  task: Task;
}

const Task = ({ task }: TaskProps) => {
  const { id, title, description, completed, deadline } = task;

  return (
    <Card className="gap-4 py-4">
      <CardHeader>
        <CardTitle className="font-normal">
          <StatusBadge
            completed={completed}
            deadline={new Date(deadline)}
            className="rounded-md px-2! py-0.5! text-xs"
          />
        </CardTitle>
        <DropdownMenu>
          <CardAction>
            <DropdownMenuTrigger asChild className="cursor-pointer">
              <IconDots className="size-5" />
            </DropdownMenuTrigger>
          </CardAction>
          <DropdownMenuContent className="w-40 px-2">
            <DropdownMenuGroup>
              {!completed && <MarkAsDoneBtn id={id} />}
              <DropdownMenuItem asChild>
                <Link href={ROUTES.editTask(id)}>
                  <IconPencil />
                  Edit
                </Link>
              </DropdownMenuItem>
            </DropdownMenuGroup>
            <DropdownMenuSeparator />
            <DropdownMenuGroup>
              <DeleteBtn id={id} />
            </DropdownMenuGroup>
          </DropdownMenuContent>
        </DropdownMenu>
      </CardHeader>
      <CardContent>
        <Link href={ROUTES.task(id)} className="space-y-2">
          <p className="font-semibold line-clamp-1">{title}</p>
          <p className="text-sm line-clamp-1 text-muted-foreground">
            {description}
          </p>
        </Link>
      </CardContent>
      <CardFooter className="mt-auto flex-col gap-3 items-start">
        <Separator />
        <div className="flex gap-2 items-center text-sm text-muted-foreground">
          <IconCalendarWeek className="size-5" />
          {format(new Date(deadline), "MMM dd HH:mm")}
        </div>
      </CardFooter>
    </Card>
  );
};

export default Task;
