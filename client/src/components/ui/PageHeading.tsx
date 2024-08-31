import { ReactNode } from "react";

interface Props {
  heading: string;
  children?: ReactNode;
}

function PageHeading({ heading, children }: Props) {
  return (
    <section className="flex justify-between items-center mb-4">
      <h1 className="text-2xl">{heading}</h1>
      {children && <span>{children}</span>}
    </section>
  );
}

export default PageHeading;
