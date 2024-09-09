import { memo } from "react";
import { Link } from "react-router-dom";
import { NavigationAuthenticationProps, NavigationProps } from ".";
import LinkButton from "../ui/LinkButton";
import { Button } from "../ui/button.tsx";
import { NavLinkProps } from "./navigationMenuItems.ts";

const DesktopNavigation = memo(({ isAuthenticated, menuItems, logout }: NavigationProps) => (
  <nav className="flex items-center gap-8">
    <ul className="space-x-8">
      {menuItems.map((link) => (
        <DesktopNavLink key={link.label} {...link} />
      ))}
    </ul>
    <AdditionalButtons isAuthenticated={isAuthenticated} logout={logout} />
  </nav>
));


const AdditionalButtons = ({ isAuthenticated, logout }: NavigationAuthenticationProps) => {
  return (
    <div className="flex gap-2">
      {
        !isAuthenticated ?
          <>
            <LinkButton to="/login" variant="outline" size="sm">
              Log in
            </LinkButton>
            <LinkButton to="/signup" size="sm">
              Sign up
            </LinkButton>
          </> :
          <>
            <Button variant="outline" onClick={logout}>Log out</Button>
            <LinkButton to="/settings">Settings</LinkButton>
          </>
      }
    </div>
  )
}

const DesktopNavLink = ({ href, label }: NavLinkProps) => (
  <Link
    to={href}
    className="text-slate-800 hover:text-purple-500 transition-colors"
  >
    {label}
  </Link>
);

export default DesktopNavigation;
