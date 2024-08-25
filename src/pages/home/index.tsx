import { Button } from "../../components/ui/button";
import PageHeading from "../../components/ui/PageHeading";

function HomePage() {
  return (
    <>
      <PageHeading heading="Welcome">
        <Button size="sm">Next</Button>
      </PageHeading>
      <p>
        Advancely is your one stop shop for everything performance and
        progression.
      </p>
    </>
  );
}

export default HomePage;
