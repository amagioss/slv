import { useEffect, useRef, type ReactNode } from 'react';
import Layout from '@theme/Layout';
import Link from '@docusaurus/Link';
import HomepageFeatures from '@site/src/components/HomepageFeatures/HomepageFeatures';
import OpenSourceSection from '@site/src/components/OpenSourceSection/OpenSourceSection';
import styles from './index.module.css';

export default function Home(): ReactNode {
  const featuresRef = useRef<HTMLDivElement>(null);
  const openSourceRef = useRef<HTMLDivElement>(null);
  const scrollContainerRef = useRef<HTMLDivElement>(null);
  const scrollIndicatorRef = useRef<HTMLButtonElement>(null);
  const featuresScrollIndicatorRef = useRef<HTMLButtonElement>(null);

  const scrollToFeatures = () => {
    if (featuresRef.current && scrollContainerRef.current) {
      featuresRef.current.scrollIntoView({ behavior: 'smooth', block: 'start' });
    }
  };

  const scrollToOpenSource = () => {
    if (openSourceRef.current && scrollContainerRef.current) {
      openSourceRef.current.scrollIntoView({ behavior: 'smooth', block: 'start' });
    }
  };

  useEffect(() => {
    // Add class to body for home page styling
    document.body.classList.add('home-page');
    
    // Function to hide navbar
    const hideNavbar = () => {
      const navbar = document.querySelector('nav.navbar') as HTMLElement;
      if (navbar) {
        navbar.style.display = 'none';
      }
    };
    
    // Hide navbar immediately
    hideNavbar();
    
    // Also try after a short delay in case navbar loads late
    const timeoutId = setTimeout(hideNavbar, 100);
    
    // Use MutationObserver to catch dynamically added navbar
    const observer = new MutationObserver(hideNavbar);
    observer.observe(document.body, {
      childList: true,
      subtree: true,
    });

    // Wait for initial animation to complete before allowing JS control
    // Animation delay (0.9s) + duration (1s) = 1.9s
    const animationCompleteTimeout = setTimeout(() => {
      if (scrollIndicatorRef.current) {
        scrollIndicatorRef.current.setAttribute('data-animation-complete', 'true');
      }
    }, 1900);

    // Intersection Observer for scroll-triggered animations
    let featuresIntersectionObserver: IntersectionObserver | null = null;
    let openSourceIntersectionObserver: IntersectionObserver | null = null;
    
    const setupIntersectionObservers = () => {
      const featuresElement = featuresRef.current;
      const openSourceElement = openSourceRef.current;
      
      if (featuresElement) {
        featuresIntersectionObserver = new IntersectionObserver(
          (entries) => {
            entries.forEach((entry) => {
              if (entry.isIntersecting) {
                entry.target.setAttribute('data-features-visible', 'true');
                // Hide landing scroll indicator when features are visible
                if (scrollIndicatorRef.current && scrollIndicatorRef.current.getAttribute('data-animation-complete') === 'true') {
                  scrollIndicatorRef.current.style.opacity = '0';
                  scrollIndicatorRef.current.style.pointerEvents = 'none';
                }
                // Show features scroll indicator when features section is visible
                if (featuresScrollIndicatorRef.current) {
                  featuresScrollIndicatorRef.current.style.opacity = '0.7';
                  featuresScrollIndicatorRef.current.style.pointerEvents = 'auto';
                }
              } else {
                // Show landing scroll indicator when features are not visible (only after animation completes)
                if (scrollIndicatorRef.current && scrollIndicatorRef.current.getAttribute('data-animation-complete') === 'true') {
                  scrollIndicatorRef.current.style.opacity = '0.7';
                  scrollIndicatorRef.current.style.pointerEvents = 'auto';
                }
                // Hide features scroll indicator when features are not visible
                if (featuresScrollIndicatorRef.current) {
                  featuresScrollIndicatorRef.current.style.opacity = '0';
                  featuresScrollIndicatorRef.current.style.pointerEvents = 'none';
                }
              }
            });
          },
          { threshold: 0.1, rootMargin: '0px' }
        );
        
        // Small delay to ensure element is rendered
        setTimeout(() => {
          if (featuresElement && featuresIntersectionObserver) {
            featuresIntersectionObserver.observe(featuresElement);
          }
        }, 100);
      }

      if (openSourceElement) {
        openSourceIntersectionObserver = new IntersectionObserver(
          (entries) => {
            entries.forEach((entry) => {
              if (entry.isIntersecting) {
                entry.target.setAttribute('data-open-source-visible', 'true');
                // Hide features scroll indicator when open source section is visible
                if (featuresScrollIndicatorRef.current) {
                  featuresScrollIndicatorRef.current.style.opacity = '0';
                  featuresScrollIndicatorRef.current.style.pointerEvents = 'none';
                }
              } else {
                // Show features scroll indicator when open source is not visible
                if (featuresScrollIndicatorRef.current && featuresRef.current) {
                  const featuresRect = featuresRef.current.getBoundingClientRect();
                  const isFeaturesVisible = featuresRect.top < window.innerHeight && featuresRect.bottom > 0;
                  if (isFeaturesVisible) {
                    featuresScrollIndicatorRef.current.style.opacity = '0.7';
                    featuresScrollIndicatorRef.current.style.pointerEvents = 'auto';
                  }
                }
              }
            });
          },
          { threshold: 0.1, rootMargin: '0px' }
        );
        
        setTimeout(() => {
          if (openSourceElement && openSourceIntersectionObserver) {
            openSourceIntersectionObserver.observe(openSourceElement);
          }
        }, 100);
      }
    };

    setupIntersectionObservers();
    
    return () => {
      clearTimeout(timeoutId);
      clearTimeout(animationCompleteTimeout);
      observer.disconnect();
      if (featuresIntersectionObserver) {
        featuresIntersectionObserver.disconnect();
      }
      if (openSourceIntersectionObserver) {
        openSourceIntersectionObserver.disconnect();
      }
      document.body.classList.remove('home-page');
      const navbar = document.querySelector('nav.navbar') as HTMLElement;
      if (navbar) {
        navbar.style.display = '';
      }
    };
  }, []);

  return (
    <Layout
      title={`${siteConfig.title}`}
      description="Secure, reliable, and scalable vault management for Kubernetes and beyond."
      noFooter={true}>
      
      {/* Hero section */}
      <HomepageHeader />

      {/* Main content */}
      <main>
        {/* Features section */}
        <HomepageFeatures />
      </main>
    </Layout>
  );
}
