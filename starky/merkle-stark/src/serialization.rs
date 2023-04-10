use std::io::Cursor;
use std::io::{Result, Write};

use plonky2::field::extension::{Extendable, FieldExtension};
use plonky2::field::types::PrimeField64;
use plonky2::fri::proof::{FriInitialTreeProof, FriProof, FriQueryRound, FriQueryStep};
use plonky2::hash::hash_types::RichField;
use plonky2::hash::merkle_proofs::MerkleProof;
use plonky2::hash::merkle_tree::MerkleCap;
use plonky2::plonk::config::{GenericConfig, GenericHashOut, Hasher};

use crate::proof::{StarkOpeningSet, StarkProof, StarkProofWithPublicInputs};

#[derive(Debug)]
pub struct Buffer(Cursor<Vec<u8>>);

impl Buffer {
    pub fn new(buffer: Vec<u8>) -> Self {
        Self(Cursor::new(buffer))
    }

    pub fn len(&self) -> usize {
        self.0.get_ref().len()
    }

    pub fn bytes(self) -> Vec<u8> {
        self.0.into_inner()
    }

    fn write_u8(&mut self, x: u8) -> Result<()> {
        self.0.write_all(&[x])
    }

    fn write_field<F: PrimeField64>(&mut self, x: F) -> Result<()> {
        self.0.write_all(&x.to_canonical_u64().to_le_bytes())
    }

    fn write_field_ext<F: RichField + Extendable<D>, const D: usize>(
        &mut self,
        x: F::Extension,
    ) -> Result<()> {
        for &a in &x.to_basefield_array() {
            self.write_field(a)?;
        }
        Ok(())
    }

    fn write_hash<F: RichField, H: Hasher<F>>(&mut self, h: H::Hash) -> Result<()> {
        self.0.write_all(&h.to_bytes())
    }

    fn write_merkle_cap<F: RichField, H: Hasher<F>>(
        &mut self,
        cap: &MerkleCap<F, H>,
    ) -> Result<()> {
        for &a in &cap.0 {
            self.write_hash::<F, H>(a)?;
        }
        Ok(())
    }

    pub fn write_field_vec<F: PrimeField64>(&mut self, v: &[F]) -> Result<()> {
        for &a in v {
            self.write_field(a)?;
        }
        Ok(())
    }

    fn write_field_ext_vec<F: RichField + Extendable<D>, const D: usize>(
        &mut self,
        v: &[F::Extension],
    ) -> Result<()> {
        for &a in v {
            self.write_field_ext::<F, D>(a)?;
        }
        Ok(())
    }

    fn write_merkle_proof<F: RichField, H: Hasher<F>>(
        &mut self,
        p: &MerkleProof<F, H>,
    ) -> Result<()> {
        let length = p.siblings.len();
        self.write_u8(
            length
                .try_into()
                .expect("Merkle proof length must fit in u8."),
        )?;
        for &h in &p.siblings {
            self.write_hash::<F, H>(h)?;
        }
        Ok(())
    }

    fn write_fri_initial_proof<
        F: RichField + Extendable<D>,
        C: GenericConfig<D, F = F>,
        const D: usize,
    >(
        &mut self,
        fitp: &FriInitialTreeProof<F, C::Hasher>,
    ) -> Result<()> {
        for (v, p) in &fitp.evals_proofs {
            self.write_field_vec(v)?;
            self.write_merkle_proof(p)?;
        }
        Ok(())
    }

    fn write_fri_query_step<
        F: RichField + Extendable<D>,
        C: GenericConfig<D, F = F>,
        const D: usize,
    >(
        &mut self,
        fqs: &FriQueryStep<F, C::Hasher, D>,
    ) -> Result<()> {
        self.write_field_ext_vec::<F, D>(&fqs.evals)?;
        self.write_merkle_proof(&fqs.merkle_proof)
    }

    fn write_fri_query_rounds<
        F: RichField + Extendable<D>,
        C: GenericConfig<D, F = F>,
        const D: usize,
    >(
        &mut self,
        fqrs: &[FriQueryRound<F, C::Hasher, D>],
    ) -> Result<()> {
        for fqr in fqrs {
            self.write_fri_initial_proof::<F, C, D>(&fqr.initial_trees_proof)?;
            for fqs in &fqr.steps {
                self.write_fri_query_step::<F, C, D>(fqs)?;
            }
        }
        Ok(())
    }

    fn write_fri_proof<F: RichField + Extendable<D>, C: GenericConfig<D, F = F>, const D: usize>(
        &mut self,
        fp: &FriProof<F, C::Hasher, D>,
    ) -> Result<()> {
        for cap in &fp.commit_phase_merkle_caps {
            self.write_merkle_cap(cap)?;
        }
        self.write_fri_query_rounds::<F, C, D>(&fp.query_round_proofs)?;
        self.write_field_ext_vec::<F, D>(&fp.final_poly.coeffs)?;
        self.write_field(fp.pow_witness)
    }

    pub fn write_stark_proof_with_public_inputs<
        F: RichField + Extendable<D>,
        C: GenericConfig<D, F = F>,
        const D: usize,
    >(
        &mut self,
        proof_with_pis: &StarkProofWithPublicInputs<F, C, D>,
    ) -> Result<()> {
        let StarkProofWithPublicInputs {
            proof,
            public_inputs,
        } = proof_with_pis;
        self.write_stark_proof(proof)?;
        self.write_field_vec(public_inputs)
    }

    pub fn write_stark_proof<
        F: RichField + Extendable<D>,
        C: GenericConfig<D, F = F>,
        const D: usize,
    >(
        &mut self,
        proof: &StarkProof<F, C, D>,
    ) -> Result<()> {
        self.write_merkle_cap(&proof.trace_cap)?;

        match &proof.permutation_zs_cap {
            Some(cap) => self.write_merkle_cap(cap)?,
            None => {}
        }

        self.write_merkle_cap(&proof.quotient_polys_cap)?;
        self.write_openings(&proof.openings)?;
        self.write_fri_proof::<F, C, D>(&proof.opening_proof)
    }

    pub fn write_openings<F: RichField + Extendable<D>, const D: usize>(
        &mut self,
        os: &StarkOpeningSet<F, D>,
    ) -> Result<()> {
        self.write_field_ext_vec::<F, D>(&os.local_values)?;
        self.write_field_ext_vec::<F, D>(&os.next_values)?;

        match &os.permutation_zs {
            Some(zs) => self.write_field_ext_vec::<F, D>(zs)?,
            None => {}
        }

        match &os.permutation_zs_next {
            Some(zs_next) => self.write_field_ext_vec::<F, D>(zs_next)?,
            None => {}
        }

        self.write_field_ext_vec::<F, D>(&os.quotient_polys)
    }
}
