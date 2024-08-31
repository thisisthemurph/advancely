import clsx from "clsx";
import { memo } from "react";
import { Link, LinkProps } from "react-router-dom";
import {
  Sheet,
  SheetClose,
  SheetContent,
  SheetHeader,
  SheetTrigger,
} from "../ui/sheet";
import Logo from "../ui/Logo";
import { Card } from "../ui/card";
import { Button } from "../ui/button";
import { NavLinkProps, NavLinks } from ".";
import { SheetCloseLinkButton as SheetLinkButton } from "../ui/LinkButton";

interface MobileNavigationProps {
  menuItems: NavLinks;
}

const MobileNavigation = memo(({ menuItems }: MobileNavigationProps) => (
  <Sheet>
    <SheetTrigger asChild>
      <Button
        variant="ghost"
        size="icon"
        className="top-3 right-3 fixed"
        aria-label="open menu"
      >
        <HamburgerIcon />
      </Button>
    </SheetTrigger>
    <SheetContent side="top" hideClose>
      <SheetHeader>
        <MobileNavLink to="/" className="mb-8">
          <Logo />
        </MobileNavLink>
      </SheetHeader>
      <MobileNavMenu items={menuItems} />
    </SheetContent>
  </Sheet>
));

function HamburgerIcon() {
  return (
    <svg
      xmlns="http://www.w3.org/2000/svg"
      fill="none"
      viewBox="0 0 24 24"
      strokeWidth={1.5}
      stroke="currentColor"
      className="size-6"
      aria-hidden="true"
    >
      <path
        strokeLinecap="round"
        strokeLinejoin="round"
        d="M3.75 6.75h16.5M3.75 12h16.5m-16.5 5.25h16.5"
      />
    </svg>
  );
}

function MobileNavMenu({ items }: { items: NavLinks }) {
  return (
    <div className="flex flex-col gap-8">
      <div className="gap-2 grid grid-cols-2">
        {items.map((link) => (
          <CardLink key={link.label} {...link} />
        ))}
      </div>
      <div className="flex justify-between gap-2">
        <SheetLinkButton to="/login" variant="outline" className="grow">
          Log in
        </SheetLinkButton>
        <SheetLinkButton to="/signup" className="grow">
          Sign up
        </SheetLinkButton>
      </div>
    </div>
  );
}

function MobileNavLink(props: LinkProps) {
  return (
    <SheetClose asChild>
      <Link {...props} />
    </SheetClose>
  );
}

function CardLink({ label, description, href, fullWidth }: NavLinkProps) {
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

export default MobileNavigation;
