import React from "react";

export default function useInfiniteScroll(
  hasMore: boolean,
  fetchData: () => void
) {
  const [loading, setLoading] = React.useState(false);
  const observer = React.useRef<IntersectionObserver | null>(null);

  React.useEffect(() => {
    if (loading) {
      setTimeout(() => {
        fetchData();

        setLoading(false);
      }, 700);
    }
  }, [loading]);

  const lastElementRef = React.useCallback(
    (node: Element | null) => {
      if (observer.current) observer.current.disconnect();
      observer.current = new IntersectionObserver((entries) => {
        if (entries[0].isIntersecting && hasMore) {
          setLoading(true);
        }
      });
      if (node) observer.current.observe(node);
    },
    [hasMore]
  );

  return [lastElementRef, loading];
}
