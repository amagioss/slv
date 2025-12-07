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
      clearTimeout(animationCompleteTimeout);
      if (featuresIntersectionObserver) {
        featuresIntersectionObserver.disconnect();
      }
      if (openSourceIntersectionObserver) {
        openSourceIntersectionObserver.disconnect();
      }
      document.body.classList.remove('home-page');
    };
  }, []);

  return (
    <Layout
      title="Home"
      description="Home page"
      noFooter>
      <div className={styles.fixedBackground}></div>
      <div className={styles.scrollContainer} ref={scrollContainerRef}>
        <section className={styles.landingSection}>
          <div className={styles.bannerSection}>
            <img 
              src="/img/slv-banner.svg" 
              alt="SLV Banner" 
              className={styles.bannerLogo}
            />
            <p className={styles.subtitle}>Securely store, share, and access secrets alongside the codebase.</p>
            <div className={styles.buttonGroup}>
              <Link
                to="/docs/quick-start"
                className={styles.getStartedButton}>
                Get Started
              </Link>
              <a
                href="https://github.com/amagioss/slv"
                className={styles.githubButton}
                target="_blank"
                rel="noopener noreferrer"
                aria-label="GitHub repository">
                <svg width="20" height="20" viewBox="0 0 24 24" fill="currentColor">
                  <path d="M12 .297c-6.63 0-12 5.373-12 12 0 5.303 3.438 9.8 8.205 11.385.6.113.82-.258.82-.577 0-.285-.01-1.04-.015-2.04-3.338.724-4.042-1.61-4.042-1.61C4.422 18.07 3.633 17.7 3.633 17.7c-1.087-.744.084-.729.084-.729 1.205.084 1.838 1.236 1.838 1.236 1.07 1.835 2.809 1.305 3.495.998.108-.776.417-1.305.76-1.605-2.665-.3-5.466-1.332-5.466-5.93 0-1.31.465-2.38 1.235-3.22-.135-.303-.54-1.523.105-3.176 0 0 1.005-.322 3.3 1.23.96-.267 1.98-.399 3-.405 1.02.006 2.04.138 3 .405 2.28-1.552 3.285-1.23 3.285-1.23.645 1.653.24 2.873.12 3.176.765.84 1.23 1.91 1.23 3.22 0 4.61-2.805 5.625-5.475 5.92.42.36.81 1.096.81 2.22 0 1.606-.015 2.896-.015 3.286 0 .315.21.69.825.57C20.565 22.092 24 17.592 24 12.297c0-6.627-5.373-12-12-12z"/>
                </svg>
                GitHub
              </a>
            </div>
          </div>
          <button 
            ref={scrollIndicatorRef}
            className={styles.scrollIndicator}
            onClick={scrollToFeatures}
            aria-label="Scroll to features">
            <span className={styles.scrollText}>Explore Features</span>
            <svg 
              className={styles.scrollArrow}
              width="24" 
              height="24" 
              viewBox="0 0 24 24" 
              fill="none" 
              stroke="currentColor" 
              strokeWidth="2" 
              strokeLinecap="round" 
              strokeLinejoin="round">
              <path d="M6 9l6 6 6-6"/>
            </svg>
          </button>
        </section>
        <section className={styles.featuresSection} ref={featuresRef}>
          <HomepageFeatures />
          <button 
            ref={featuresScrollIndicatorRef}
            className={styles.featuresScrollIndicator}
            onClick={scrollToOpenSource}
            aria-label="Scroll to open source section">
            <svg 
              className={styles.scrollArrow}
              width="24" 
              height="24" 
              viewBox="0 0 24 24" 
              fill="none" 
              stroke="currentColor" 
              strokeWidth="2" 
              strokeLinecap="round" 
              strokeLinejoin="round">
              <path d="M6 9l6 6 6-6"/>
            </svg>
          </button>
        </section>
        <section className={styles.openSourceSection} ref={openSourceRef}>
          <OpenSourceSection />
        </section>
      </div>
    </Layout>
  );
}
