import React, { useEffect, useRef, useState } from 'react';
import Link from '@docusaurus/Link';
import styles from './OpenSourceSection.module.css';
import { GITHUB_STATS } from './githubStats';

interface GitHubStats {
  stars: number;
  contributors: number;
}

export default function OpenSourceSection(): React.JSX.Element {
  const [isVisible, setIsVisible] = useState(false);
  const [stats, setStats] = useState<GitHubStats>({ 
    stars: GITHUB_STATS.stars, 
    contributors: GITHUB_STATS.contributors 
  });
  const sectionRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    // Try to fetch GitHub stats from API
    const fetchStats = async () => {
      try {
        // Fetch repository info for stars
        const repoResponse = await fetch('https://api.github.com/repos/amagioss/slv');
        
        if (repoResponse.status === 403) {
          // Rate limited - use build-time values
          console.log('GitHub API rate limited, using build-time values');
          return;
        }
        
        if (repoResponse.ok) {
          const repoData = await repoResponse.json();
          const stars = repoData.stargazers_count || 0;
          setStats(prev => ({ ...prev, stars }));
        }

        // Fetch contributors count
        let contributorsCount = 0;
        let page = 1;
        let hasMore = true;

        while (hasMore && page <= 10) {
          const contributorsResponse = await fetch(
            `https://api.github.com/repos/amagioss/slv/contributors?per_page=100&page=${page}`
          );
          
          if (contributorsResponse.status === 403) {
            // Rate limited - use build-time values
            console.log('GitHub API rate limited, using build-time values');
            break;
          }
          
          if (contributorsResponse.ok) {
            const contributorsData = await contributorsResponse.json();
            
            if (contributorsData.length === 0) {
              hasMore = false;
            } else {
              contributorsCount += contributorsData.length;
              const linkHeader = contributorsResponse.headers.get('Link');
              if (!linkHeader || !linkHeader.includes('rel="next"')) {
                hasMore = false;
              } else {
                page++;
              }
            }
          } else {
            hasMore = false;
          }
        }

        if (contributorsCount > 0) {
          setStats(prev => ({ ...prev, contributors: contributorsCount }));
        }
      } catch (error) {
        // Network error or other issue - use build-time values
        console.log('Failed to fetch GitHub stats, using build-time values:', error);
      }
    };

    fetchStats();
  }, []);

  useEffect(() => {
    // Set up Intersection Observer
    const observer = new IntersectionObserver(
      (entries) => {
        entries.forEach((entry) => {
          if (entry.isIntersecting) {
            setIsVisible(true);
            entry.target.setAttribute('data-oss-visible', 'true');
          }
        });
      },
      { threshold: 0.2, rootMargin: '0px' }
    );

    if (sectionRef.current) {
      observer.observe(sectionRef.current);
    }

    return () => {
      if (sectionRef.current) {
        observer.unobserve(sectionRef.current);
      }
    };
  }, []);

  return (
    <div ref={sectionRef} className={styles.openSourceSection}>
      <div className={styles.content}>
        <h2 className={styles.title}>100% Open Source</h2>
        <p className={styles.description}>
          Built by the community, for the community. SLV is free, open source, and always will be.
        </p>

        <div className={styles.statsContainer}>
          <div className={styles.statCard}>
            <div className={styles.statIcon}>⭐</div>
            <div className={styles.statValue}>
              {isVisible && stats.stars > 0 ? (
                <AnimatedCounter key={`stars-${stats.stars}`} target={stats.stars} delay={1200} />
              ) : (
                stats.stars > 0 ? stats.stars.toLocaleString() : '0'
              )}
            </div>
            <div className={styles.statLabel}>GitHub Stars</div>
          </div>

          <div className={styles.statCard}>
            <div className={styles.statIcon}>👥</div>
            <div className={styles.statValue}>
              {isVisible && stats.contributors > 0 ? (
                <AnimatedCounter key={`contributors-${stats.contributors}`} target={stats.contributors} delay={1200} />
              ) : (
                stats.contributors > 0 ? stats.contributors.toLocaleString() : '0'
              )}
            </div>
            <div className={styles.statLabel}>Contributors</div>
          </div>
        </div>

        <div className={styles.actionButtons}>
          <Link
            to="/docs/quick-start"
            className={styles.actionButton}>
            Quick Start
          </Link>
          <a
            href="https://slv.sh/docs/contributing"
            className={styles.actionButton}
            target="_blank"
            rel="noopener noreferrer">
            Contribution Guide
          </a>
          <a
            href="#"
            className={styles.actionButton}>
            Blog
          </a>
        </div>

        <div className={styles.repoLink}>
          <a
            href="https://github.com/amagioss/slv"
            target="_blank"
            rel="noopener noreferrer"
            className={styles.repoButton}>
            <svg width="20" height="20" viewBox="0 0 24 24" fill="currentColor">
              <path d="M12 .297c-6.63 0-12 5.373-12 12 0 5.303 3.438 9.8 8.205 11.385.6.113.82-.258.82-.577 0-.285-.01-1.04-.015-2.04-3.338.724-4.042-1.61-4.042-1.61C4.422 18.07 3.633 17.7 3.633 17.7c-1.087-.744.084-.729.084-.729 1.205.084 1.838 1.236 1.838 1.236 1.07 1.835 2.809 1.305 3.495.998.108-.776.417-1.305.76-1.605-2.665-.3-5.466-1.332-5.466-5.93 0-1.31.465-2.38 1.235-3.22-.135-.303-.54-1.523.105-3.176 0 0 1.005-.322 3.3 1.23.96-.267 1.98-.399 3-.405 1.02.006 2.04.138 3 .405 2.28-1.552 3.285-1.23 3.285-1.23.645 1.653.24 2.873.12 3.176.765.84 1.23 1.91 1.23 3.22 0 4.61-2.805 5.625-5.475 5.92.42.36.81 1.096.81 2.22 0 1.606-.015 2.896-.015 3.286 0 .315.21.69.825.57C20.565 22.092 24 17.592 24 12.297c0-6.627-5.373-12-12-12z"/>
            </svg>
            View on GitHub
          </a>
        </div>
      </div>
    </div>
  );
}

// Animated counter component
function AnimatedCounter({ target, delay = 0 }: { target: number; delay?: number }): React.JSX.Element {
  const [count, setCount] = useState(0);
  const [hasStarted, setHasStarted] = useState(false);

  useEffect(() => {
    // Reset state when target changes
    setCount(0);
    setHasStarted(false);

    if (target === 0) return;

    // Wait for the delay before starting the animation
    const startTimer = setTimeout(() => {
      setHasStarted(true);
    }, delay);

    return () => clearTimeout(startTimer);
  }, [target, delay]);

  useEffect(() => {
    if (!hasStarted || target === 0) return;

    const duration = 2000; // 2 seconds
    const steps = 60;
    const increment = target / steps;
    const stepDuration = duration / steps;

    let current = 0;
    const timer = setInterval(() => {
      current += increment;
      if (current >= target) {
        setCount(target);
        clearInterval(timer);
      } else {
        setCount(Math.floor(current));
      }
    }, stepDuration);

    return () => clearInterval(timer);
  }, [target, hasStarted]);

  return <>{count.toLocaleString()}</>;
}

