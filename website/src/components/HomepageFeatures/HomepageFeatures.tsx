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
    title: 'Kubernetes Native',
    Svg: require('@site/static/img/kubernetes.svg').default,
    description: (
      <>
        SLV integrates seamlessly into Kubernetes as a native custom resource (CRD), eliminating the need for external vault servers or complex sidecars.
      </>
    ),
  },
  {
    title: '"Secrets As Code" Ready',
    Svg: require('@site/static/img/git-branch.svg').default,
    description: (
      <>
        Safely commit encrypted secrets into your Git repositories without risk. SLV enables true GitOps workflows, treating secrets as part of your infrastructure code.
      </>
    ),
  },
  {
    title: 'Complete SDLC Integration',
    Svg: require('@site/static/img/code.svg').default,
    description: (
      <>
        From local CLI development to GitHub Actions to Kubernetes clusters, SLV secures your secrets throughout the entire software development lifecycle.
      </>
    ),
  },
];

function Feature({ title, Svg, description }: FeatureItem) {
  return (
    <div className={clsx('col col--4', styles.feature)}>
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
      <div className="container">
        <div className="row">
          {FeatureList.map((props, idx) => (
            <Feature key={idx} {...props} />
          ))}
        </div>
      </div>
    </section>
  );
}
