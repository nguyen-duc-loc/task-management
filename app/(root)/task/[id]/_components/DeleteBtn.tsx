"use client";

import { deleteTask } from "@/api/actions/task";
import { Button } from "@/components/ui/button";
import Spinner from "@/components/ui/spinner";
import ROUTES from "@/constants/routes";
import { IconTrash } from "@tabler/icons-react";
import { useRouter } from "next/navigation";
import React from "react";
import { toast } from "sonner";

interface DeleteBtnProps {
  id: string;
}

const DeleteBtn = ({ id }: DeleteBtnProps) => {
  const [deletingTask, startDeletingTask] = React.useTransition();
  const router = useRouter();

  const handleMarkTaskAsDone = () => {
    startDeletingTask(async () => {
      const result = await deleteTask(id);
      if (!result.success) {
        toast.error(`Failed to delete task. Try again later.`);
        return;
      }

      toast.success("Task has been deleted!");
      router.push(ROUTES.dashboard);
    });
  };

  return (
    <Button
      variant="destructive"
      disabled={deletingTask}
      onClick={handleMarkTaskAsDone}
    >
      {deletingTask ? <Spinner /> : <IconTrash />}
      Delete
    </Button>
  );
};

export default DeleteBtn;
