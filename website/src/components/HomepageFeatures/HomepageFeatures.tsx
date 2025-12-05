import React from 'react';
import styles from './HomepageFeatures.module.css';
import clsx from 'clsx';

type FeatureItem = {
  title: string;
  Svg: React.ComponentType<React.ComponentProps<'svg'>>;
  description: JSX.Element;
};

const FeatureList: FeatureItem[] = [
  {
    title: 'Complete SDLC Integration',
    Svg: require('@site/static/img/code.svg').default,
    description: (
      <>
        From local development to CI builds and deployment (cloud or on-prem), SLV secures your secrets across the entire software development lifecycle.
      </>
    ),
  },
  {
    title: '"Secrets As Code" Ready',
    Svg: require('@site/static/img/git-branch.svg').default,
    description: (
      <>
        Securely commit SLV-encrypted secrets into Git repositories with confidence. SLV enables true GitOps by managing secrets as code, similar to configs.
      </>
    ),
  },
  {
    title: 'Kubernetes Native',
    Svg: require('@site/static/img/kubernetes.svg').default,
    description: (
      <>
        SLV seamlessly integrates with Kubernetes via a custom resource automatically converted to native secrets, eliminating external vault servers or complex sidecars.
      </>
    ),
  },
  {
    title: 'Quantum Resistant',
    Svg: require('@site/static/img/qc.svg').default,
    description: (
      <>
        Built with future-proof cryptography that resists quantum computing attacks, ensuring your secrets remain secure even as technology evolves.
      </>
    ),
  },
];

function Feature({ title, Svg, description }: FeatureItem) {
  return (
    <div className={styles.feature}>
      <div className="text--center">
        <Svg className={styles.featureSvg} role="img" />
      </div>
      <div className="text--center padding-horiz--md">
        <h3>{title}</h3>
        <p>{description}</p>
      </div>
    </div>
  );
}

export default function HomepageFeatures(): JSX.Element {
  return (
    <section className={styles.features}>
      <div className={styles.featuresGrid}>
        {FeatureList.map((props, idx) => (
          <Feature key={idx} {...props} />
        ))}
      </div>
    </section>
  );
}
