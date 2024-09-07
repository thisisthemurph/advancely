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
import { Card } from "../ui/card";
import { Button } from "../ui/button";
import { NavLinkProps, NavLinks } from ".";
import { SheetCloseLinkButton as SheetLinkButton } from "../ui/LinkButton";
import Logo from "../Logo.tsx";

interface MobileNavigationProps {
  menuItems: NavLinks;
  isAuthenticated: boolean;
  logout: () => Promise<void>;
}

const MobileNavigation = memo(({ menuItems, isAuthenticated, logout }: MobileNavigationProps) => {
  return (
    <Sheet>
      <SheetTrigger asChild>
        <Button
          variant="ghost"
          size="icon"
          aria-label="open menu"
        >
          <HamburgerIcon />
        </Button>
      </SheetTrigger>
      <SheetContent side="top" hideClose>
        <SheetHeader>
          <MobileNavLink to="/" className="group flex justify-center items-center gap-2 mb-8">
            <Logo size="md" className="group-hover:grayscale group-hover:translate-x-[5.5rem] transition-all"/>
            <span className="text-slate-600 italic font-semibold group-hover:-translate-x-12 transition-all">Advancely</span>
          </MobileNavLink>
        </SheetHeader>
        <MobileNavMenu items={menuItems} isAuthenticated={isAuthenticated} logout={logout} />
      </SheetContent>
    </Sheet>
  )
});

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

function MobileNavMenu({ items, isAuthenticated, logout }: { items: NavLinks, isAuthenticated: boolean, logout: () => Promise<void> }) {
  return (
    <div className="flex flex-col gap-8">
      <div className="gap-2 grid grid-cols-2">
        {items.map((link) => (
          <CardLink key={link.label} {...link} />
        ))}
      </div>
      <AdditionalButtons isAuthenticated={isAuthenticated} logout={logout} />
    </div>
  );
}

const AdditionalButtons = ({isAuthenticated, logout}: {isAuthenticated: boolean, logout: () => Promise<void>}) => {
  return (
    <div className="flex justify-between gap-2">
      {!isAuthenticated && (
        <>
          <SheetLinkButton to="/login" variant="outline" className="grow">
            Log in
          </SheetLinkButton>
          <SheetLinkButton to="/signup" className="grow">
            Sign up
          </SheetLinkButton>
        </>
      )}
      {isAuthenticated && (
        <>
          <SheetClose asChild>
            <Button variant="outline" className="grow" onClick={logout}>Log out</Button>
          </SheetClose>
          <SheetLinkButton to="/settings" className="grow">Settings</SheetLinkButton>
        </>
      )}
    </div>
  )
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
