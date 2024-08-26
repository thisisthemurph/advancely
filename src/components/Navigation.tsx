import {
  Sheet,
  SheetClose,
  SheetContent,
  SheetHeader,
  SheetTrigger,
} from "./ui/sheet";
import { Button } from "./ui/button";
import LinkButton, {
  SheetCloseLinkButton as SheetLinkButton,
} from "./ui/LinkButton";
import { Card } from "./ui/card";
import { Link, LinkProps } from "react-router-dom";
import clsx from "clsx";
import Logo from "./ui/Logo";

import { useMediaQuery } from "usehooks-ts";

const cardMenuItems: CardLinkProps[] = [
  {
    label: "Home",
    href: "/",
    description: "There's no place like home!",
  },
  {
    label: "About",
    href: "/",
    description: "Find out a little of who we are and what we can do for you.",
  },
  {
    label: "Pricing",
    href: "/",
    description:
      "Take a look at our competative pricing models and see what best suits your needs.",
    fullWidth: true,
  },
];

function Navigation() {
  const isDesktop = useMediaQuery("(min-width: 640px)");

  if (isDesktop) {
    return <DesktopNavigation />;
  }

  return <MobileNavigation />;
}

function DesktopNavigation() {
  return (
    <nav className="flex items-center gap-8">
      <ul className="space-x-8">
        {cardMenuItems.map((link, index) => (
          <Link
            key={index}
            to={link.href}
            className="text-slate-800 hover:text-purple-500 transition-colors"
          >
            {link.label}
          </Link>
        ))}
      </ul>
      <div className="flex gap-2">
        <LinkButton to="/login" variant="outline" size="sm">
          Log in
        </LinkButton>
        <LinkButton to="/signup" size="sm">
          Sign up
        </LinkButton>
      </div>
    </nav>
  );
}

function MobileNavigation() {
  return (
    <Sheet>
      <SheetTrigger asChild>
        <Button variant="ghost" size="icon" className="top-3 right-3 fixed">
          <svg
            xmlns="http://www.w3.org/2000/svg"
            fill="none"
            viewBox="0 0 24 24"
            strokeWidth={1.5}
            stroke="currentColor"
            className="size-6"
          >
            <path
              strokeLinecap="round"
              strokeLinejoin="round"
              d="M3.75 6.75h16.5M3.75 12h16.5m-16.5 5.25h16.5"
            />
          </svg>
        </Button>
      </SheetTrigger>
      <SheetContent side="top" hideClose>
        <SheetHeader>
          <SheetCloseLink to="/" className="mb-8">
            <Logo />
          </SheetCloseLink>
        </SheetHeader>
        <div className="flex flex-col gap-8">
          <div className="gap-2 grid grid-cols-2">
            {cardMenuItems.map((link, index) => (
              <CardLink {...link} key={index} />
            ))}
          </div>
          <div className="flex justify-between gap-4">
            <SheetLinkButton to="/login" variant="outline" className="grow">
              Log in
            </SheetLinkButton>
            <SheetLinkButton to="/signup" className="grow">
              Sign up
            </SheetLinkButton>
          </div>
        </div>
      </SheetContent>
    </Sheet>
  );
}

interface CardLinkProps {
  label: string;
  description: string;
  href: string;
  fullWidth?: boolean;
}

function SheetCloseLink(props: LinkProps) {
  return (
    <SheetClose asChild>
      <Link {...props} />
    </SheetClose>
  );
}

function CardLink({ label, description, href, fullWidth }: CardLinkProps) {
  return (
    <Card
      className={clsx(
        "hover:bg-gray-50 transition-colors",
        fullWidth && "col-span-2"
      )}
    >
      <SheetClose asChild>
        <Link to={href} className="flex flex-col gap-2 p-4 w-full">
          <p className="text-xl">{label}</p>
          <p className="text-slate-600 text-sm">{description}</p>
        </Link>
      </SheetClose>
    </Card>
  );
}

export default Navigation;
