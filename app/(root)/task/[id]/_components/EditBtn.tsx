import { Button } from "@/components/ui/button";
import ROUTES from "@/constants/routes";
import { IconEdit } from "@tabler/icons-react";
import Link from "next/link";
import React from "react";

interface EditTaskBtnProps {
  id: string;
}

const EditTaskBtn = ({ id }: EditTaskBtnProps) => {
  return (
    <Button asChild>
      <Link href={ROUTES.editTask(id)}>
        <IconEdit />
        Edit
      </Link>
    </Button>
  );
};

export default EditTaskBtn;
