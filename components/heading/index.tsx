import React from "react";

interface HeadingProps {
  heading: string;
  Icon?: React.ComponentType<{ className: string }>;
}

const Heading = ({ heading, Icon }: HeadingProps) => {
  return (
    <h1 className="mb-12 flex items-center gap-2 text-2xl font-bold">
      {Icon && <Icon className="size-7" />}
      {heading}
    </h1>
  );
};

export default Heading;
